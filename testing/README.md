# Testing

Docker image with all the required tools to run the ci/release actions

# Build
`docker build -f testing/Dockerfile -t openintegration/testing .`

# Release
`docker push openintegration/testing`