#!/bin/bash
cd /home/jian/reddit-crawler
docker compose -f /home/jian/reddit-crawler/docker-compose.yml up
git add .
git commit -m "Updated"
GIT_SSH_COMMAND="ssh -i ~/.ssh/id_ed25519 -F /dev/null" git push
