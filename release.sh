cd ../
mkdir clockify-tracker-release
cd clockify-tracker-release
mkdir bin

cd ../clockify-time-tracker

# build for all platforms
GOOS=darwin GOARCH=amd64 go build -o ../clockify-tracker-release/bin/clockify-tracker-mac-intel
GOOS=darwin GOARCH=arm64 go build -o ../clockify-tracker-release/bin/clockify-tracker-mac-arm
GOOS=linux GOARCH=amd64 go build -o ../clockify-tracker-release/bin/clockify-tracker-linux
GOOS=windows GOARCH=amd64 go build -o ../clockify-tracker-release/bin/clockify-tracker.exe

# Copy documentation
cp .env.example ../clockify-tracker-release/
cp README.md ../clockify-tracker-release/
cp INSTALL.md ../clockify-tracker-release/

# Create INSTALL.md in the release folder (copy from the artifact I created earlier)

# Go back to workspace and create zip
cd ..
zip -r clockify-tracker-release.zip clockify-tracker-release/
