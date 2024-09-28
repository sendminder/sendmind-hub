docker build --platform linux/amd64 -t sendmind-hub:latest .
docker tag sendmind-hub:latest seninder/sendmind-hub:latest
docker push seninder/sendmind-hub:latest
