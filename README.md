# quick-response

a go qr code generator deployed on vercel with dockerfile.vercel.

testing the dockerfile on vercel feature:
https://vercel.com/blog/dockerfile-on-vercel

vercel auto-detected the dockerfile preset during project setup:

![app-preset](screenshots/app-preset.png)

## endpoints

- `GET /` - web form to generate qr codes
- `GET /qr?text=hello&size=256` - returns a png
- `GET /health` - health check

## deploy

```bash
vercel deploy
```

or push to github and vercel picks it up.

## what this tests

- building a go binary inside a container on vercel
- storing and serving container images
- env vars like $port
- query params and binary responses
- cold starts and redeploys
