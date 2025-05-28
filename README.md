Task Manager is a web-based application for managing tasks and visualizing task-related data. It features a user-friendly interface to view and complete tasks, along with data visualizations for task pass/fail rates, turn-around times, monthly open/close counts, and task type distributions. The application runs two servers: a Python Flask API for data processing and a Go server for serving the frontend and proxying API requests. A launcher script with a GUI allows starting and stopping both servers with a single button.
Features

Task Management: View up to 20 tasks, mark as "Completed" or "Cannot Complete" with comments, and track completion timestamps.
Data Visualizations:
Pass/Fail: Pie chart showing the latest pass/fail counts.
Turn Around Time: Line chart of average turn-around time (in days) for the last 5 months.
Open/Close Monthly: Line chart of total open and closed tasks per month for the last 5 months.
Task Types: Pie chart showing the distribution of task types (e.g., "Bug Fix: 10").
Tasks Completed This Week: Stat showing tasks completed in the current week.


Single-Button Control: Start and stop both servers using a Tkinter GUI.
Cross-Platform: Runs on Windows, Linux, and macOS, with a packaged executable option.

Dependencies

Python 3.8+:
Flask (pip install flask)
openpyxl (pip install openpyxl)
python-dateutil (pip install python-dateutil)
PyInstaller (pip install pyinstaller) for packaging
Tkinter (usually included; install via sudo apt-get install python3-tk on Ubuntu or brew install python-tk on macOS if missing)


Go 1.16+: For compiling the Go server.
Browser: Any modern browser (Chrome, Firefox, etc.) for the web interface.

Ports

Python API: Runs on http://localhost:5000
Go Server: Runs on http://localhost:8080 (access the web interface here)

File Structure
task_manager/

├── api.py                # Python Flask API for data processing

├── main.go               # Go server for frontend and API proxy

├── static/
│   └── index.html        # Frontend HTML with Chart.js visualizations

├── tasks.xlsx            # Excel file with task data (generated)

├── launcher.py           # Tkinter GUI to start/stop servers

├── server                # Compiled Go binary (server.exe on Windows)

└── server.log            # Log file for server output


Prerequisites

Install Python 3.8+ from python.org.
Install Go 1.16+ from golang.org.
Ensure git is installed to clone the repository.

Setup

Clone the Repository:
git clone https://github.com/vkrause09/task_manager.git
cd task_manager


Install Python Dependencies:
pip install flask openpyxl python-dateutil pyinstaller

Compile the Go Server:
# Linux/Mac
GOOS=linux GOARCH=amd64 go build -o server main.go
# Windows
GOOS=windows GOARCH=amd64 go build -o server.exe main.go

This creates server (or server.exe) in the task_manager/ directory.


Running Locally

Start the Servers:
python launcher.py


A GUI opens with "Start Servers" and "Stop Servers" buttons.
Click "Start Servers":
Copies tasks.xlsx to the working directory.
Starts Python API (http://localhost:5000) and Go server (http://localhost:8080).
Opens http://localhost:8080 in your default browser.


Logs are written to server.log.


Access the Application:

Open http://localhost:8080 (if not opened automatically).
Features:
Tasks To Do: View and complete tasks (e.g., 4 tasks).
Data:
"Tasks Completed This Week: 3" (top).
2x2 charts: Pass/Fail, Turn Around Time, Open/Close Monthly, Task Types.


Header shows current time (e.g., "8:24:00 PM CDT").




Stop the Servers:

In the GUI, click "Stop Servers" to gracefully shut down both servers.
Alternatively, close the GUI window to stop servers.



Packaging as a Single Executable
To distribute the application as a single executable (no Python/Go required on the target machine):

Bundle with PyInstaller:
pyinstaller --onefile \
    --add-data "api.py:." \
    --add-data "server:." \
    --add-data "static/index.html:static" \
    --add-data "tasks.xlsx:." \
    --name task_manager \
    launcher.py


On Windows, use ; instead of ::pyinstaller --onefile ^
    --add-data "api.py;." ^
    --add-data "server.exe;." ^
    --add-data "static/index.html;static" ^
    --add-data "tasks.xlsx;." ^
    --name task_manager ^
    launcher.py


Output: dist/task_manager (or dist/task_manager.exe).


Run the Executable:
./dist/task_manager  # Linux/Mac
dist\task_manager.exe  # Windows


Opens the GUI.
Click "Start Servers" to launch the application.
Click "Stop Servers" to shut down.



Troubleshooting

Empty Tasks ({"all":[],"completed":[]}):
Cause: tasks.xlsx is missing or misplaced.
Fix:
Re-run python create_mock_tasks.py and ensure tasks.xlsx is in task_manager/.
Check server.log for errors like "File tasks.xlsx does not exist".
Verify tasks.xlsx contains all sheets (AllTasks, Types, etc.).


Test:curl http://localhost:5000/tasks

Expect: {"all":[{"name":"Write report",...},...],"completed":[...]}


Servers Don’t Start:
Check server.log for errors.
Ensure ports 5000 and 8080 are free:netstat -an | grep 5000
netstat -an | grep 8080


Verify server (or server.exe) is executable (chmod +x server on Linux/Mac).


Tkinter Missing:
Install:sudo apt-get install python3-tk  # Ubuntu
brew install python-tk  # macOS




PyInstaller Issues:
Ensure all --add-data paths match file locations.
Check dist/task_manager/_internal/ for api.py, server, tasks.xlsx.



Contributing
Contributions are welcome! Please:

Fork the repository.
Create a feature branch (git checkout -b feature/YourFeature).
Commit changes (git commit -m "Add YourFeature").
Push to the branch (git push origin feature/YourFeature).
Open a Pull Request.

License
MIT License
Copyright (c) 2025 Victor Krause
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

AI Training Policy

The use of this repository’s code, documentation, or data for training artificial intelligence models, machine learning systems, or similar technologies is prohibited without explicit written permission from the repository owner. If you wish to use this project for AI training purposes, please contact victorkrause9@gmail.com to discuss licensing terms, which may include a paid agreement. This restriction does not affect other uses permitted under the MIT License, such as personal or commercial use, modification, or distribution, provided the license terms are followed.
