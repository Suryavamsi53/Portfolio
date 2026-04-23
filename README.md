# 🚀 SVV Portfolio - Production Ready

A premium, cloud-native portfolio with an interactive admin dashboard and global visitor tracking.

## 🛠️ Features
- **Interactive UI**: Glassmorphism design with Lucide icons.
- **Admin Dashboard**: Gmail-style inbox for managing contact transmissions.
- **Live World Map**: Interactive visualization of visitor locations using Leaflet.js.
- **Secure Authentication**: Go-based backend with JWT and Bcrypt hashing.
- **Demo Mode**: Full functionality on static servers (GitHub Pages/S3) with mock fallbacks.

## 📦 Deployment Options

### 1. GitHub Pages (Static Hosting)
This is the easiest way to host your portfolio.
- **Login**: Use `admin` / `admin123`.
- **Contact Form**: Works in "Demo Mode" (simulates success).
- **Dashboard**: Works in "Demo Mode" using data from `messages.json`.

**Steps to Deploy:**
1. Push this folder to a GitHub Repository.
2. Go to **Settings > Pages**.
3. Select the `main` branch as the source.
4. Your site will be live at `https://<username>.github.io/<repo>/`.

---

### 2. Full Server Deployment (Go Backend)
For a fully functional dashboard where you can receive and delete messages.
- **Requirements**: A server with Go installed (EC2, Heroku, DigitalOcean, etc.)

**Steps to Run:**
1. Set the following Environment Variables on your server:
   - `PORT=8080` (or any port)
   - `JWT_SECRET=your-secure-random-key`
2. Run the server:
   ```bash
   go run main.go
   ```
3. Access your portfolio at `http://your-server-ip:8080`.

---

## 🔒 Security Notes
- The `JWT_SECRET` should be a long, random string in production.
- Default Admin credentials (`admin`/`admin123`) should be changed in `main.go` before a real server deployment.

## 🎨 Technology Stack
- **Frontend**: HTML5, Vanilla CSS, JS (ES6+)
- **Backend**: Go (Golang)
- **Maps**: Leaflet.js + OpenStreetMap (Nominatim)
- **Icons**: Lucide Icons
