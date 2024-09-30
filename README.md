<p align="center">
  <img src=".github/1.webp" alt="Blum Automation Screenshot" width="300"/>
</p>

# Blum Automation Software

This software automates various activities in the Blum app, including tasks, gaming, and daily claims. The app interacts with the Blum platform and performs these actions automatically based on the configuration you set.

## Features

- **Task Automation**: Automatically complete daily tasks on Blum.
- **Gaming Automation**: Play games and claim rewards based on available play passes.
- **Daily Claims**: Automatically claim daily check-in rewards and farming rewards.

## Setup

### Prerequisites

- Go (version 1.18+)
- Telegram's Blum WebApp must be run with web inspection enabled.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/malvere/Blumfield.git
   cd Blumfield
   ```

2. Install Go dependencies:
   ```bash
   go mod download
   ```

3. Configure your settings:
   The configuration file `config.yaml` must be filled with your WebApp Init data and settings for automation. The init data is retrieved by running Blum in the Telegram web app with web inspection enabled.

### Obtaining WebApp Init Data

1. Open the Blum app in your **Telegram Web** app.
2. Enable web inspection (right-click the page and select **Inspect** or press `F12`).
3. In the console tab of the inspector, type:
   ```javascript
   copy(Telegram.WebApp.initData)
   ```
4. Paste the copied output into the `configs/config.yaml` file under the `auth` section:
   ```yaml
   auth:
     WebAppInit: "your_copied_webapp_init_data"
   ```

### Configuration

In the `config.yaml` file, you can set various parameters for the automation process:

```yaml
auth:
  WebAppInit: ""
  tokenFile: "tokens.json"
  folder: "config"

settings:
  daemon: false         # Run perpetually (restart every 8 hrs)
  randomAgent: true     # Randomise User-Agent
  delay: 3              # Delay between tasks (in seconds)
  tasks: true           # Enable or disable task automation
  farming: true         # Enable or disable farming claims
  gaming: true          # Enable or disable gaming automation
```

### Running the Software

Once the configuration is set, you can run the software using:

```bash
./build/blumfield
```

### Logging

The software logs activity to the console and provides details on the actions being taken, such as claiming rewards, completing tasks, and farming.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## Contributions

Feel free to submit issues or pull requests if you have ideas for improvements!

---

This `README.md` provides all necessary setup instructions, including how to obtain and configure the `WebAppInit` data and run the project. You can modify it based on any additional project details.
