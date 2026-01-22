# URL Shortener — Docker & Railway Deploy

This repository runs a small Go-based URL shortener that serves a single page UI and an API on port `8080`.

This README shows how to build a Docker image locally, push it to a registry (Docker Hub), and deploy to Railway.

## Files added
- `Dockerfile` — multi-stage build producing a small static Linux binary and copying `index.html` into the final image.
- `.dockerignore` — excludes unnecessary files from build context.

## Build locally

Open PowerShell in the repository root (where `Dockerfile` is) and run:

```powershell
docker build -t your-dockerhub-username/urlshortener:latest .
```

Replace `your-dockerhub-username` with your Docker Hub account.

## Run locally to test

```powershell
docker run -p 8080:8080 your-dockerhub-username/urlshortener:latest
```

Then open `http://localhost:8080` in your browser.

## Push to Docker Hub

1. Log in:

```powershell
docker login
```

2. Push the image:

```powershell
docker push your-dockerhub-username/urlshortener:latest
```

Now the image is available publicly (or privately if your repo is private).

## Deploy on Railway — two options

Option A — Deploy from GitHub (recommended):

1. Commit and push your repository to GitHub.
2. Go to https://railway.app, create a project.
3. Choose to deploy from GitHub and connect your repo.
4. Railway detects the `Dockerfile` and will build the image automatically. If asked, set the service port to `8080`.

Option B — Deploy from a local Docker image using Railway CLI:

1. Install Railway CLI (follow instructions on Railway docs) or run `npm i -g @railway/cli` if you have Node.js.
2. Login and initialize a project:

```powershell
railway login
railway init
```

3. Deploy the current folder (Railway will use the `Dockerfile`):

```powershell
railway up
```

Railway will build and run the container and provide a public URL.

## Notes and troubleshooting
- `index.html` is read at runtime from the container root; the `Dockerfile` copies it into the final image.
- If you get `permission denied` when running the binary, ensure the image was built with `CGO_ENABLED=0` (the provided `Dockerfile` does this).
- If Railway reports build errors, open the build logs to see missing files or Go build output. Common fixes: ensure `go.mod` is present and all source files are in the repo.

If you want, I can:

- Build the Docker image here and run it to verify (I can show the exact commands to run locally).
- Push to Docker Hub (you'll need to provide credentials or push from your machine).
- Walk through connecting the GitHub repo to Railway and setting environment variables, DNS, or custom domain.
