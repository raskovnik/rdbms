package templates

import (
    "html/template"
    "net/http"
)

var indexHTML = `<!DOCTYPE html>
<html>
<head>
    <title>TODO List - RDBMS Demo</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 600px;
            margin: 50px auto;
            padding: 20px;
        }
        h1 { color: #333; }
        .add-form {
            display: flex;
            gap: 10px;
            margin-bottom: 20px;
        }
        input[type="text"] {
            flex: 1;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
        }
        button {
            padding: 10px 20px;
            background: #007bff;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }
        button:hover { background: #0056b3; }
        .todo-item {
            display: flex;
            align-items: center;
            padding: 10px;
            border: 1px solid #ddd;
            margin-bottom: 10px;
            border-radius: 4px;
        }
        .todo-item.completed { opacity: 0.6; }
        .todo-item input[type="checkbox"] {
            margin-right: 10px;
        }
        .todo-text {
            flex: 1;
        }
        .todo-item.completed .todo-text {
            text-decoration: line-through;
        }
        .delete-btn {
            background: #dc3545;
            padding: 5px 10px;
        }
        .delete-btn:hover { background: #c82333; }
        .error {
            color: #dc3545;
            margin: 10px 0;
        }
    </style>
</head>
<body>
    <h1>TODO List</h1>
    <p><em>Powered by custom RDBMS</em></p>
    
    <div class="add-form">
        <input type="text" id="taskInput" placeholder="Enter a new task..." />
        <button onclick="addTodo()">Add</button>
    </div>
    
    <div id="error" class="error"></div>
    <div id="todoList"></div>
    
    <script>
        loadTodos();
        
        async function loadTodos() {
            try {
                const response = await fetch('/todos');
                const todos = await response.json();
                
                const list = document.getElementById('todoList');
                list.innerHTML = '';
                
                if (!todos || todos.length === 0) {
                    list.innerHTML = '<p>No todos yet. Add one above!</p>';
                    return;
                }
                
                todos.forEach(todo => {
                    const item = document.createElement('div');
                    item.className = 'todo-item' + (todo.completed ? ' completed' : '');
                    
                    item.innerHTML = ` + "`" + `
                        <input type="checkbox" 
                               ${todo.completed ? 'checked' : ''} 
                               onchange="toggleTodo(${todo.id}, this.checked)">
                        <span class="todo-text">${escapeHtml(todo.task)}</span>
                        <button class="delete-btn" onclick="deleteTodo(${todo.id})">Delete</button>
                    ` + "`" + `;
                    
                    list.appendChild(item);
                });
            } catch (error) {
                showError('Failed to load todos: ' + error.message);
            }
        }
        
        async function addTodo() {
            const input = document.getElementById('taskInput');
            const task = input.value.trim();
            
            if (!task) {
                showError('Task cannot be empty');
                return;
            }
            
            try {
                const response = await fetch('/todos', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ task })
                });
                
                if (!response.ok) throw new Error('Failed to create todo');
                
                input.value = '';
                loadTodos();
                clearError();
            } catch (error) {
                showError(error.message);
            }
        }
        
        async function toggleTodo(id, completed) {
            try {
                const response = await fetch(` + "`" + `/todos/${id}` + "`" + `, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ completed: completed ? 1 : 0 })
                });
                
                if (!response.ok) throw new Error('Failed to update todo');
                
                loadTodos();
            } catch (error) {
                showError(error.message);
            }
        }
        
        async function deleteTodo(id) {
            if (!confirm('Delete this todo?')) return;
            
            try {
                const response = await fetch(` + "`" + `/todos/${id}` + "`" + `, {
                    method: 'DELETE'
                });
                
                if (!response.ok) throw new Error('Failed to delete todo');
                
                loadTodos();
            } catch (error) {
                showError(error.message);
            }
        }
        
        function showError(msg) {
            document.getElementById('error').textContent = msg;
        }
        
        function clearError() {
            document.getElementById('error').textContent = '';
        }
        
        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }
        
        document.getElementById('taskInput').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') addTodo();
        });
    </script>
</body>
</html>`

func ServeIndex(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.New("index").Parse(indexHTML))
    w.Header().Set("Content-Type", "text/html")
    tmpl.Execute(w, nil)
}