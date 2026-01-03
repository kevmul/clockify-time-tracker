# Clockify Time Tracker - Installation Guide

## Quick Start

### 1. Get Your Clockify API Key

1. Log in to [Clockify](https://clockify.me)
2. Go to Settings → Profile
3. Scroll down to "API" section
4. Copy your API key

### 2. Choose Your Platform

**macOS (Intel):**

```bash
# Copy the binary
cp bin/clockify-tracker-mac-intel /usr/local/bin/clockify-tracker
chmod +x /usr/local/bin/clockify-tracker
```

**macOS (Apple Silicon - M1/M2/M3):**

```bash
# Copy the binary
cp bin/clockify-tracker-mac-arm /usr/local/bin/clockify-tracker
chmod +x /usr/local/bin/clockify-tracker
```

**Linux:**

```bash
# Copy the binary
sudo cp bin/clockify-tracker-linux /usr/local/bin/clockify-tracker
sudo chmod +x /usr/local/bin/clockify-tracker
```

**Windows:**

1. Copy `bin/clockify-tracker.exe` to `C:\Program Files\ClockifyTracker\`
2. Add `C:\Program Files\ClockifyTracker\` to your PATH:
   - Search for "Environment Variables" in Windows
   - Edit "Path" under System Variables
   - Add new entry: `C:\Program Files\ClockifyTracker`
   - Restart your terminal

### 3. Configure Your API Key

Create a `.env` file in your home directory:

**macOS/Linux:**

```bash
# Create config directory
mkdir -p ~/.config/clockify-tracker

# Copy the example file
cp .env.example ~/.config/clockify-tracker/.env

# Edit with your API key
nano ~/.config/clockify-tracker/.env
# or
vim ~/.config/clockify-tracker/.env
# or use any text editor
```

**Windows:**

```cmd
# Create config directory
mkdir %USERPROFILE%\.config\clockify-tracker

# Copy the example file
copy .env.example %USERPROFILE%\.config\clockify-tracker\.env

# Edit with your API key using Notepad
notepad %USERPROFILE%\.config\clockify-tracker\.env
```

Add your API key to the file:

```
CLOCKIFY_API_KEY=your_api_key_here
```

### 4. Run It!

Just type from anywhere:

```bash
clockify-tracker
```

## Troubleshooting

### "command not found: clockify-tracker"

**macOS/Linux:**

- Make sure `/usr/local/bin` is in your PATH
- Check: `echo $PATH`
- If not, add to your `~/.bashrc` or `~/.zshrc`:
  ```bash
  export PATH="/usr/local/bin:$PATH"
  ```

**Windows:**

- Make sure you added the folder to your PATH
- Restart your terminal/PowerShell after changing PATH

### "CLOCKIFY_API_KEY not set"

The tool looks for `.env` in these locations (in order):

1. Current directory (`./.env`)
2. Home config directory (`~/.config/clockify-tracker/.env`)

Make sure you've created the `.env` file in one of these locations.

### Permission Denied (macOS/Linux)

```bash
chmod +x /usr/local/bin/clockify-tracker
```

## Alternative: Run from Current Directory

If you don't want to add to PATH, you can run from the directory:

```bash
# Set up in a folder
mkdir ~/clockify-tracker
cd ~/clockify-tracker
cp /path/to/clockify-tracker-mac-intel ./clockify-tracker
cp /path/to/.env.example ./.env
# Edit .env with your API key
nano .env

# Run it
./clockify-tracker
```

## Uninstall

**macOS/Linux:**

```bash
rm /usr/local/bin/clockify-tracker
rm -rf ~/.config/clockify-tracker
```

**Windows:**

```cmd
del "C:\Program Files\ClockifyTracker\clockify-tracker.exe"
rmdir /s "C:\Program Files\ClockifyTracker"
# Remove from PATH in Environment Variables
```
