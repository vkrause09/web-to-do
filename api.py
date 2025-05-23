from flask import Flask, jsonify
from openpyxl import load_workbook
from datetime import datetime

app = Flask(__name__)
file_path = 'tasks.xlsx'

def read_tasks():
    wb = load_workbook(file_path)
    all_tasks_sheet = wb['AllTasks']
    completed_tasks_sheet = wb['CompletedTasks']

    tasks = {
        "all": [],
        "completed": []
    }

    def parse_tasks(sheet):
        task_list = []
        for row in sheet.iter_rows(min_row=2, values_only=True):
            task = {
                "name": row[0],
                "priority": row[1],
                "date_added": row[2]
            }
            task_list.append(task)
        return task_list

    tasks["all"] = sorted(
        parse_tasks(all_tasks_sheet),
        key=lambda x: (x['priority'], x['date_added'])
    )
    tasks["completed"] = parse_tasks(completed_tasks_sheet)

    return tasks

def mark_task_completed(task_name):
    wb = load_workbook(file_path)
    all_tasks_sheet = wb['AllTasks']
    completed_tasks_sheet = wb['CompletedTasks']

    for row in all_tasks_sheet.iter_rows(min_row=2, values_only=False):
        if row[0].value == task_name:
            completed_row = [cell.value for cell in row]
            completed_row[2] = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
            completed_tasks_sheet.append(completed_row)
            all_tasks_sheet.delete_rows(row[0].row)
            break

    wb.save(file_path)

def get_recent_pass_fail():
    wb = load_workbook(file_path)
    pass_fail_sheet = wb['Pass Fail']

    most_recent_row = None
    for row in pass_fail_sheet.iter_rows(min_row=2, values_only=True):
        if most_recent_row is None or row[0] > most_recent_row[0]:
            most_recent_row = row

    if most_recent_row:
        return {
            "date": most_recent_row[0],
            "pass": most_recent_row[1],
            "fail": most_recent_row[2]
        }
    return {}

@app.route('/tasks', methods=['GET'])
def get_tasks():
    tasks = read_tasks()
    return jsonify(tasks)

@app.route('/tasks/complete', methods=['POST'])
def complete_task():
    task_name = request.json.get('name')
    mark_task_completed(task_name)
    return '', 204

@app.route('/pass_fail', methods=['GET'])
def pass_fail():
    data = get_recent_pass_fail()
    return jsonify(data)

if __name__ == '__main__':
    app.run(port=5000)
