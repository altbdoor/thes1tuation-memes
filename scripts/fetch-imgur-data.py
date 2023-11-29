#!/usr/bin/env python

from urllib.request import Request, urlopen
import os
import json
import random
import shutil
from typing import TypedDict
from itertools import groupby
from datetime import datetime
from zoneinfo import ZoneInfo

client_id = os.getenv("IMGUR_CLIENT_ID")
album_hash = "yzKq60n"
current_dir = os.path.dirname(__file__)
imgur_json = os.path.join(current_dir, "../_data/imgur-parsed.json")
static_copy_json = os.path.join(current_dir, "../assets/imgur.json")

user_agent = (
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) "
    f"Chrome/79.0.3945.{random.randint(0, 9999)} Safari/537.{random.randint(0, 99)}"
)

if client_id is None:
    raise Exception("please provide IMGUR_CLIENT_ID")
else:
    print("(i) calling somewhat official imgur API")
    req = Request(
        f"https://api.imgur.com/post/v1/albums/{album_hash}?include=media,tags,account",
        None,
        {"Authorization": f"Client-ID {client_id}", "User-Agent": user_agent},
    )

    with urlopen(req) as res:
        data = json.loads(res.read())

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
RawImage = TypedDict(
    "RawImage",
    {
        "id": str,
        "created_at": str,
        "mime_type": str,
        "width": int,
        "height": int,
        "size": int,
        "url": str,
    },
)

remapped_data: list[Image] = []

for img in data.get("media", []):
    raw_image: RawImage = {key: img.get(key) for key in RawImage.__annotations__.keys()}

    img_datetime = datetime.fromisoformat(raw_image["created_at"])
    img_datetime = img_datetime.replace(tzinfo=ZoneInfo("UTC")).astimezone(
        tz=ZoneInfo("Asia/Kuala_Lumpur")
    )

    thumbnail = (
        raw_image["url"]
        .replace(raw_image["id"], raw_image["id"] + "b")
        .replace(".gif", ".jpg")
        .replace(".jpeg", ".jpg")
    )

    remapped_img: Image = {
        "id": raw_image["id"],
        "datetime": int(img_datetime.timestamp()),
        "type": raw_image["mime_type"],
        "width": raw_image["width"],
        "height": raw_image["height"],
        "size": raw_image["size"],
        "link": raw_image["url"].replace(".jpeg", ".jpg"),
        "thumbnail": thumbnail,
        "timeDisplay": img_datetime.strftime("%d %b, %I:%M %p"),
        "groupBy": img_datetime.strftime("%B %Y"),
    }

    remapped_data.append(remapped_img)

with open(imgur_json, "w", newline="") as fp:
    remapped_data = sorted(remapped_data, key=lambda x: x["datetime"], reverse=True)
    groups = groupby(remapped_data, key=lambda x: x["groupBy"])

    group_data = [{"name": key, "items": list(items)} for key, items in groups]
    json.dump(group_data, fp, indent=2)
    print(f"grouped {len(remapped_data)} images into {len(group_data)} months")

# give a static copy to web
shutil.copy(imgur_json, static_copy_json)
