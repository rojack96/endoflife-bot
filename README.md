## 📘 **EndOfLife Discord Bot**

### 🧠 Overview

**EndOfLife Bot** is a Discord bot that integrates with the [endoflife.date](https://endoflife.date) public API to help you track software lifecycle information directly in your Discord server.

With this bot, you can quickly check:

* When a product version reaches **End of Life (EOL)**
* Which releases are currently **supported**
* Full **release histories** and **LTS versions**

---

### 💬 Commands

| Command                                 | Description                                                                                                             |
| --------------------------------------- | ----------------------------------------------------------------------------------------------------------------------- |
| `/help`                                 | Shows all available commands and their usage                                                                            |
| `/product-list [page]`                  | Displays a paginated list of all products available from [endoflife.date](https://endoflife.date). Default page is `1`. |
| `/product-lts <product>`                | Displays Long Term Support (LTS) information for a specific product.                                                    |
| `/product-info <product> [page]`        | Shows detailed release information for a product, paginated by version. Default page is `1`.                            |
| `/product-releases <product> <release>` | Displays detailed information about a specific product release, or use `latest` to view the most recent one.            |

---

### 🧩 Example Usage

```
/product-list 2
```

→ Shows the second page of available products.

```
/product-lts ubuntu
```

→ Shows LTS information for Ubuntu.

```
/product-info nodejs 3
```

→ Lists all Node.js releases on page 3.

```
/product-releases nodejs latest
```

→ Shows detailed info about the latest Node.js release.

---

### 🛠️ Technologies Used

* **Golang** 🦦 — backend logic and API handling
* [**discordgo**](https://github.com/bwmarrin/discordgo) — Discord API SDK
* [**endoflife.date API**](https://endoflife.date/docs/api/v1/) — public data source for software lifecycle info

---

### 🧾 License

This project is open-source and available under the **MIT License**.