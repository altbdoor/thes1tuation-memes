#!/usr/bin/env python

from urllib.request import Request, urlopen
import os
import json
from typing import TypedDict
from itertools import groupby
from datetime import datetime, timezone

client_id = os.getenv("IMGUR_CLIENT_ID")
album_hash = "yzKq60n"
current_dir = os.path.dirname(__file__)
cache_json = os.path.join(current_dir, "cache.json")
imgur_json = os.path.join(current_dir, "../data/imgur.json")

if client_id is None:
    raise Exception("please provide client id")

req = Request(
    f"https://api.imgur.com/3/album/{album_hash}/images",
    None,
    {"Authorization": f"Client-ID {client_id}"},
)
with urlopen(req) as res, open(cache_json, "w") as fp:
    data = json.loads(res.read())
    fp.write(json.dumps(data, indent=4))

# with open(cache_json, "r") as fp:
#     data = json.load(fp)

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
    remapped_data.append(remapped_img)

remapped_data = sorted(remapped_data, key=lambda x: x["datetime"], reverse=True)

with open(imgur_json, "w") as fp:
    groups = groupby(
        remapped_data,
        key=lambda x: datetime.fromtimestamp(x["datetime"], timezone.utc).strftime(
            "%b %Y"
        ),
    )

    group_data = [{"name": key, "items": list(items)} for key, items in groups]
    json.dump(group_data, fp, indent=4)
