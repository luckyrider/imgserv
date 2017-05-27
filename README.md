# Image server

[![Build Status](https://travis-ci.org/luckyrider/imgserv.svg?branch=master)](https://travis-ci.org/luckyrider/imgserv)

## Run

```bash
docker run --rm -v $(pwd):/src centurylink/golang-builder
docker-compose up
```

## Description

Implement an HTTP server with the following routes:

- `/generate/png/{width_px}/{height_px}` - returns a black `image/png` that's
of the specified width and height
- `/generate/jpg/{width_px}/{height_px}` - returns a black `image/jpeg` that's
of the specified width and height
- `/stats` - returns `application/json` with the number of images generated
since the server was started and the average width and height of those images.

For the `/stats` route, the json should look like this:

```json
{
  "num_images": 1,
  "average_width_px": 123,
  "average_height_px": 456
}
```
