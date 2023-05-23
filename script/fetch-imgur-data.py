#!/usr/bin/env python

from urllib.request import Request, urlopen
import os
import json
import random
from typing import TypedDict
from itertools import groupby
from datetime import datetime
from zoneinfo import ZoneInfo

client_id = os.getenv("IMGUR_CLIENT_ID")
album_hash = "yzKq60n"
current_dir = os.path.dirname(__file__)
cache_json = os.path.join(current_dir, "imgur-cache.json")
imgur_json = os.path.join(current_dir, "../data/imgur-parsed.json")

# bypass imgur rate limit by hosting the json elsewhere
bypass_ratelimit_url = os.getenv("IMGUR_BYPASS_RATELIMIT_URL")

user_agent = (
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) "
    f"Chrome/79.0.3945.{random.randint(0, 9999)} Safari/537.{random.randint(0, 99)}"
)

if bypass_ratelimit_url is not None:
    print("(i) bypassing rate limits on imgur with a third party URL")

    req = Request(
        bypass_ratelimit_url,
        None,
        {"User-Agent": user_agent}
    )
    with urlopen(req) as res, open(cache_json, "w") as fp:
        data = json.loads(res.read())
        json.dump(data, fp, indent=4)
elif client_id is not None:
    req = Request(
        f"https://api.imgur.com/3/album/{album_hash}/images",
        None,
        {
            "Authorization": f"Client-ID {client_id}",
            "User-Agent": user_agent
        },
    )
    with urlopen(req) as res, open(cache_json, "w") as fp:
        data = json.loads(res.read())
        json.dump(data, fp, indent=4)
else:
    if os.path.exists(cache_json):
        print("(i) found cached response, using cache")
        with open(cache_json, "r") as fp:
            data = json.load(fp)
    else:
        raise Exception("please provide IMGUR_CLIENT_ID")

Image = TypedDict(
    "Image",
    {
        "id": str,
        "datetime": int,
        "type": str,
        "width": int,
        "height": int,
        "size": int,
        "link": str,
        "thumbnail": str,
        "timeDisplay": str,
        "groupBy": str,
    },
)

remapped_data: list[Image] = []
valid_keys = (
    "id",
    "datetime",
    "type",
    "width",
    "height",
    "size",
    "link",
)
for img in data.get("data", []):
    remapped_img: Image = {key: img.get(key) for key in valid_keys}
    remapped_img["thumbnail"] = (
        remapped_img["link"]
        .replace(remapped_img["id"], remapped_img["id"] + "b")
        .replace(".gif", ".jpg")
    )

    img_datetime = datetime.fromtimestamp(
        remapped_img["datetime"], ZoneInfo("Asia/Kuala_Lumpur")
    )
    remapped_img["timeDisplay"] = img_datetime.strftime("%d %b, %I:%M %p")
    remapped_img["groupBy"] = img_datetime.strftime("%B %Y")

    remapped_data.append(remapped_img)

with open(imgur_json, "w") as fp:
    remapped_data = sorted(remapped_data, key=lambda x: x["datetime"], reverse=True)
    groups = groupby(remapped_data, key=lambda x: x["groupBy"])

    group_data = [{"name": key, "items": list(items)} for key, items in groups]
    json.dump(group_data, fp, indent=4)
    print(f"grouped {len(remapped_data)} images into {len(group_data)} months")
