---
title: imgur shadowban
code: true
description: |
  Some memes are haram.
---

<div class="alert alert-primary">

tldr; memes are safe in [here](../imgur/){: .alert-link}, you can find the backups here:

1. [Images (61.0 MB)](https://f000.backblazeb2.com/file/thes1tuation-memes/magomet.zip)
1. [Tags in YAML](https://github.com/altbdoor/thes1tuation-memes/blob/master/_data/imgur-tags.yml)
1. [Raw JSON data]({{ site.baseurl }}/assets/imgur.json)

</div>

### First failure, possible flooding

At about <time datetime="2023-05-23T15:17:34+08:00">May 23, 2023, 3:17 PM</time>,
I tried updating the website with some tags, but the build failed.
The HTTP error was 429, which was commonly known as imgur's response
to high traffic.

I could not be spamming imgur servers, I rarely update the site even.
But then again, I remember the [imgur NSFW fiasco][1],
which _might_ have triggered some people to start archiving from GitHub CI servers.

I tried the usual bypassing by faking the `User-Agent`, but imgur servers
replied with the same HTTP 429. That led me to believe that it is probably
an IP range block, and GitHub CI servers are blocked.

Defeated, I generated the `imgur.json` data file, hosted it on a third party
website, and get the GitHub CI servers to fetch from there instead. This worked
for once, and I thought things would improve after a week or so.

### Second failure, album is 404

At about <time datetime="2023-05-26T17:53:26+08:00">May 26, 2023, 5:53 PM</time>,
I tried updating the site again. This time, the error was... odd. It is a HTTP 404.

I tried opening the [album link](https://imgur.com/a/yzKq60n), and it works for me.
I tried replaying the API call on my local machine, to get all the album images,
and it is still HTTP 404.

```console
$ curl -I -H 'Authorization: Client-ID xxx' 'https://api.imgur.com/3/album/yzKq60n/images'
HTTP/1.1 404 Not Found
Connection: keep-alive
Content-Length: 4323
content-type: text/html; charset=utf-8
```

I was stumped. My album is there, but it is not there. I told myself,
it is only one image that is going to be updated anyways. I will wait it out
for a week or so.

For what its worth, the [imgur status page](https://status.imgur.com/) has been
"fixing" an issue since May 12.

> **Ongoing intermittent issues loading content**
>
> We are working to address ongoing issues causing degraded performance across all platforms causing errors on apps, blank screens to appear, and generally wonky behavior.
>
> &mdash; <https://status.imgur.com/incidents/0g0nq22f38yh>
{: .blockquote.ps-3.border-start.border-4.border-success}

### Third failure, album is 404 for friends

After a couple of days, I hung out in Twitch chat with a couple of friends.
I was informed that they could no longer access the album, and they see a 404 page.
Once again, works on my local machine. Curious, I opened a private tab in Firefox,
and voila, HTTP 404.

I did not receive any notice about how or why the album is returning HTTP 404,
from imgur. I think it _might_ be that one NSFW meme, but I could not be sure.
At <time datetime="2023-05-28T11:50:23+08:00">May 28, 2023, 11:50 AM</time>,
I pushed some new Python code that fetches from a different API.

```diff
+ f"https://api.imgur.com/post/v1/albums/{album_hash}?include=media,tags,account"
- f"https://api.imgur.com/3/album/{album_hash}/images"
```

This uses an officially undocumented API, which appears to be used by
[major scraping libraries](https://github.com/mikf/gallery-dl/blob/v1.25.5/gallery_dl/extractor/imgur.py#L413).
This fixes the API part, and I had to make some adjustments to the data models.

But the album is still 404, so.

### Maybe it is time to panic

I did panic a little when I heard about the [imgur NSFW fiasco][1]. I did make
_some_ preparations in case I get forced to close shop, and evaluated some
alternatives. In the end, I concluded that the images should be fine, and I would
backup the images regularly, in expectance of closing shop.

With how things are going right now, I will be publishing all the links to whatever
backup I have. Things are a little scattered, but this is not exactly a
world class project to begin with. I just rolled with whatever worked for me.

If you want to use the data, some notes:

1. The images are a direct ZIP download from imgur (e.g., <https://imgur.com/a/yzKq60n/zip>),
   so the order of the images (like the album view), is messed up. You will
   probably have to cross reference the `imgur.json` to order it properly.
1. The image ZIP does not have the thumbnails.
1. The tags are in YAML because I was too lazy to write JSON. The keys are the
   imgur ID, and values are text tags.
1. The `imgur.json` has all the basic data. To note a few, the imgur ID, created timestamp,
   size, file size, and others.

I will still be maintaining this website from time to time. If imgur really forces
me to close shop, I will re-setup things on a different host. Keeping the site
up is the least I could do.

Good luck, and have fun.

[1]: https://old.reddit.com/r/DataHoarder/comments/12sbch3/imgur_is_updating_their_tos_on_may_15_2023_all/
