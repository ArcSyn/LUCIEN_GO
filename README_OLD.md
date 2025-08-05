🧠 Lucien CLI



\*\*Lucien CLI\*\* is a mystical, modular, local-first assistant that runs in your terminal.



Built with \[Bubble Tea](https://github.com/charmbracelet/bubbletea), it combines cyberpunk UI, real-time widgets, and a plugin-based architecture for spell-like control over your machine.



> 🧙 \*“Lucien awaits...”\*



---



\## ✨ Features



\- 🧾 \*\*Command input engine\*\* with text prompt

\- 🌈 \*\*Theme loader\*\* using TOML configs

\- 🕒 \*\*Real-time clock + weather widgets\*\*

\- 🔌 \*\*Plugin system\*\* using manifest.json

\- 🔮 \*\*Mystical banner system\*\*

\- 📦 \*\*Modular layout via BubbleTea + Lipgloss\*\*



---



\## 🧩 Project Structure



LucienCLI/

├── main.go # Core app logic

├── plugins/ # External plugin folders (JSON + Python)

│ └── examplePlugin/

│ ├── manifest.json

│ └── run.py

├── themes/

│ └── default.toml # Theme config

└── README.md



yaml

Copy

Edit



---



\## 📐 Theme Example (`default.toml`)



```toml

name = "Void Blue"

bg\_color = "235"

fg\_color = "252"

accent\_color = "39"

📡 Weather API Used

Lucien uses Open-Meteo for fetching weather data.

Location coordinates are auto-resolved using the city name.



🧪 Example Plugin (manifest.json)

json

Copy

Edit

{

&nbsp; "name": "hello",

&nbsp; "description": "Greets the user.",

&nbsp; "command": "hello",

&nbsp; "exec": "python run.py"

}

🧠 Philosophy

Lucien isn’t just a CLI.



It’s an occult UI layer built to empower creators, developers, and visionaries with a terminal interface that feels like invoking spells.



🚧 Roadmap

&nbsp;🔥 lucien install <plugin> support



&nbsp;🧠 Claude/GPT memory agent integration



&nbsp;🪄 Animated startup banner + sound



&nbsp;🌐 Network daemon luciend (multi-device)



&nbsp;📂 Spell manager + docked plugin tabs



🪙 Author

Built by ArcSyn



Follow progress, plugins, and prototypes at:

📦 github.com/ArcSyn/LucienCLI



🧵 Banner Preview

less

Copy

Edit

&nbsp;+-++-++-++-++-++-+

&nbsp;|L||U||C||I||E||N|

&nbsp;+-++-++-++-++-++-+

&nbsp;+-++-++-++-++-++-+

&nbsp;|A||w||a||i||t||s|

&nbsp;+-++-++-++-++-++-+

🔓 License

MIT



yaml

Copy

Edit



