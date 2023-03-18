#!/usr/bin/env python

import os
import json

current_dir = os.path.dirname(__file__)
cache_json = os.path.join(current_dir, "imgur-cache.json")
imgur_tags = os.path.join(current_dir, "../data/imgur-tags.yml")

with open(cache_json, "r") as fp:
    data = json.load(fp)

minimal_data: list[dict[str, str]] = [
    {"id": datum.get("id"), "datetime": datum.get("datetime")} for datum in data["data"]
]
minimal_data = sorted(minimal_data, key=lambda x: x["datetime"], reverse=False)
sorted_ids = [datum.get("id", "") for datum in minimal_data]

sorted_tags: list[str] = []
with open(imgur_tags, "r") as fp:
    raw_tag_lines = fp.readlines()

for image_id in sorted_ids:
    matched_raw_tag_line = list(
        filter(lambda x: x.startswith(f"{image_id}: ["), raw_tag_lines)
    )
    if len(matched_raw_tag_line) == 0:
        print(f"unable to find tags for {image_id}, patching as []")
        matched_raw_tag_line = [f"{image_id}: []\n"]

    sorted_tags.append(matched_raw_tag_line[0])

with open(imgur_tags, "w", newline="\n") as fp:
    fp.writelines(sorted_tags)
