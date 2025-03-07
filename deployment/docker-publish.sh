#!/bin/bash

echo "Building Docker image for embedder..."

# Hardcoded paths and values
docker build -t gen-mind/cognix-embedder:latest -f ../src/ai/embedder/Dockerfile ../src/ai/embedder

# Check if the build succeeded
if [ $? -ne 0 ]; then
  echo "Failed to build embedder. Exiting..."
  exit 1
fi

# Push the Docker image
echo "Pushing Docker image for embedder..."
docker push gen-mind/cognix-embedder:latest

echo "Docker image publishing complete for embedder."



#chmod +x docker-publish.sh
#docker login
#./docker-publish.sh beta



## Find the group name for your user
#id -gn
#
## Use the group name to set the correct permissions
#sudo chown -R $(whoami):staff /Users/gp/Developer/cognix/data
#
## Set correct permissions
#sudo chmod -R 755 /Users/gp/Developer/cognix/data/
