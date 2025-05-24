import tkinter as tk
import subprocess
import webbrowser
import os
import signal
import shutil
import logging
from threading import Thread

# Set up logging
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('server.log'),
        logging.StreamHandler()
    ]
)

class ServerLauncher:
    def __init__(self, root):
        self.root = root
        self.root.title("Task Manager Launcher")
        self.root.geometry("300x150")
        self.python_proc = None
        self.go_proc = None
        self.running = False

        # GUI elements
        self.status_label = tk.Label(root, text="Servers: Stopped", font=("Arial", 12))
        self.status_label.pack(pady=10)

        self.start_button = tk.Button(
            root, text="Start Servers", command=self.start_servers, width=20
        )
        self.start_button.pack(pady=5)

        self.stop_button = tk.Button(
            root, text="Stop Servers", command=self.stop_servers, width=20, state="disabled"
        )
        self.stop_button.pack(pady=5)

        # Ensure tasks.xlsx is in the right place
        self.ensure_tasks_xlsx()

    def ensure_tasks_xlsx(self):
        """Copy tasks.xlsx to the current directory if it doesn't exist."""
        try:
            tasks_src = os.path.join(os.path.dirname(__file__), 'tasks.xlsx')
            tasks_dest = os.path.join(os.getcwd(), 'tasks.xlsx')
            if not os.path.exists(tasks_dest) and os.path.exists(tasks_src):
                shutil.copy(tasks_src, tasks_dest)
                logging.info(f"Copied tasks.xlsx to {tasks_dest}")
            elif not os.path.exists(tasks_dest):
                logging.error("tasks.xlsx not found in source directory")
            else:
                logging.debug("tasks.xlsx already exists in destination")
        except Exception as e:
            logging.error(f"Error copying tasks.xlsx: {e}")

    def start_servers(self):
        """Start both Python and Go servers."""
        if self.running:
            logging.info("Servers already running")
            return

        try:
            # Start Python server (api.py)
            python_cmd = [os.sys.executable, "api.py"]
            self.python_proc = subprocess.Popen(
                python_cmd,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                text=True,
                cwd=os.getcwd()
            )

            # Start Go server (server binary)
            go_cmd = ["./server"]
            if os.name == 'nt':  # Windows
                go_cmd = ["server.exe"]
            self.go_proc = subprocess.Popen(
                go_cmd,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                text=True,
                cwd=os.getcwd()
            )

            self.running = True
            self.status_label.config(text="Servers: Running")
            self.start_button.config(state="disabled")
            self.stop_button.config(state="normal")
            logging.info("Started Python and Go servers")

            # Open browser
            webbrowser.open("http://localhost:8080")

            # Log server output in background
            Thread(target=self.log_output, args=(self.python_proc, "Python")).start()
            Thread(target=self.log_output, args=(self.go_proc, "Go")).start()

        except Exception as e:
            logging.error(f"Error starting servers: {e}")
            self.stop_servers()

    def log_output(self, proc, name):
        """Log output from a subprocess."""
        while proc.poll() is None:
            line = proc.stdout.readline().strip()
            if line:
                logging.debug(f"{name} server: {line}")

    def stop_servers(self):
        """Stop both Python and Go servers."""
        if not self.running:
            logging.info("Servers already stopped")
            return

        try:
            # Terminate Python server
            if self.python_proc and self.python_proc.poll() is None:
                self.python_proc.terminate()
                try:
                    self.python_proc.wait(timeout=5)
                except subprocess.TimeoutExpired:
                    self.python_proc.kill()
                logging.info("Stopped Python server")

            # Terminate Go server
            if self.go_proc and self.go_proc.poll() is None:
                self.go_proc.terminate()
                try:
                    self.go_proc.wait(timeout=5)
                except subprocess.TimeoutExpired:
                    self.go_proc.kill()
                logging.info("Stopped Go server")

            self.running = False
            self.status_label.config(text="Servers: Stopped")
            self.start_button.config(state="normal")
            self.stop_button.config(state="disabled")

        except Exception as e:
            logging.error(f"Error stopping servers: {e}")

    def on_closing(self):
        """Handle window close event."""
        self.stop_servers()
        self.root.destroy()

if __name__ == "__main__":
    root = tk.Tk()
    app = ServerLauncher(root)
    root.protocol("WM_DELETE_WINDOW", app.on_closing)
    root.mainloop()
