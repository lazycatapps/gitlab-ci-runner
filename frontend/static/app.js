const API_BASE = '/api';

// Load runners on page load
document.addEventListener('DOMContentLoaded', () => {
    loadRunners();
});

// Handle form submission
document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const name = document.getElementById('name').value;
    const url = document.getElementById('url').value;
    const token = document.getElementById('token').value;

    try {
        const response = await fetch(`${API_BASE}/runners/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name, url, token }),
        });

        const data = await response.json();
        console.log(data);

        if (data.success) {
            showMessage('Runner 注册成功！', 'success');
            document.getElementById('registerForm').reset();
            loadRunners();
        } else {
            showMessage(`注册失败: ${data.message || response.statusText}`, 'error');
        }
    } catch (error) {
        showMessage(`注册失败: ${error.message}`, 'error');
    }
});

// Load runners from API
async function loadRunners() {
    const loading = document.getElementById('loading');
    const runnersTable = document.getElementById('runnersTable');
    const emptyState = document.getElementById('emptyState');
    const runnersBody = document.getElementById('runnersBody');

    loading.style.display = 'block';
    runnersTable.style.display = 'none';
    emptyState.style.display = 'none';

    try {
        const response = await fetch(`${API_BASE}/runners`);
        const runners = await response.json();

        loading.style.display = 'none';

        if (runners && runners.length > 0) {
            runnersTable.style.display = 'table';
            runnersBody.innerHTML = runners.map(runner => `
                <tr>
                    <td>${escapeHtml(runner.name)}</td>
                    <td>${escapeHtml(runner.url)}</td>
                    <td><span class="token-display">${escapeHtml(runner.token)}</span></td>
                    <td>${getStatusBadge(runner.status)}</td>
                    <td>
                        <div class="btn-group">
                            <button class="btn btn-info btn-small" onclick="viewLogs('${escapeHtml(runner.name)}')">查看日志</button>
                            <button class="btn btn-warning btn-small" onclick="restartRunner('${escapeHtml(runner.name)}')">重启</button>
                            <button class="btn btn-danger btn-small" onclick="deleteRunner('${escapeHtml(runner.name)}', '${escapeHtml(runner.token)}')">删除</button>
                        </div>
                    </td>
                </tr>
            `).join('');
        } else {
            emptyState.style.display = 'block';
        }
    } catch (error) {
        loading.style.display = 'none';
        showMessage(`加载失败: ${error.message}`, 'error');
    }
}

// Restart individual runner
async function restartRunner(name) {
    if (!confirm(`确定要重启 Runner "${name}" 吗？`)) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/runners/restart`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name }),
        });

        const data = await response.json();

        if (data.success) {
            showMessage(`Runner "${name}" 重启成功！`, 'success');
            // Reload runners after a short delay to get updated status
            setTimeout(() => loadRunners(), 1000);
        } else {
            showMessage(`重启失败: ${data.message || response.statusText}`, 'error');
        }
    } catch (error) {
        showMessage(`重启失败: ${error.message}`, 'error');
    }
}

// Store current runner name for log refresh
let currentLogRunnerName = null;

// View runner logs
async function viewLogs(name) {
    const modal = document.getElementById('logModal');
    const title = document.getElementById('logModalTitle');
    const content = document.getElementById('logContent');

    // Store runner name for refresh
    currentLogRunnerName = name;

    // Show modal
    modal.style.display = 'block';
    title.textContent = `Runner "${name}" 的日志`;
    content.textContent = '加载中...';

    try {
        const response = await fetch(`${API_BASE}/runners/logs?name=${encodeURIComponent(name)}`);
        const data = await response.json();

        if (data.success) {
            content.textContent = data.logs || '暂无日志';
        } else {
            content.textContent = `获取日志失败: ${data.message || response.statusText}`;
        }
    } catch (error) {
        content.textContent = `获取日志失败: ${error.message}`;
    }
}

// Refresh logs for current runner
async function refreshLogs() {
    if (!currentLogRunnerName) {
        return;
    }

    const content = document.getElementById('logContent');
    content.textContent = '刷新中...';

    try {
        const response = await fetch(`${API_BASE}/runners/logs?name=${encodeURIComponent(currentLogRunnerName)}`);
        const data = await response.json();

        if (data.success) {
            content.textContent = data.logs || '暂无日志';
        } else {
            content.textContent = `获取日志失败: ${data.message || response.statusText}`;
        }
    } catch (error) {
        content.textContent = `获取日志失败: ${error.message}`;
    }
}

// Close log modal
function closeLogModal() {
    const modal = document.getElementById('logModal');
    modal.style.display = 'none';
    currentLogRunnerName = null;
}

// Close modal when clicking outside of it
window.onclick = function(event) {
    const modal = document.getElementById('logModal');
    if (event.target === modal) {
        closeLogModal();
    }
}

// Delete runner
async function deleteRunner(name, token) {
    if (!confirm('确定要删除这个 Runner 吗？')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/runners/delete`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name, token }),
        });

        const data = await response.json();

        if (data.success) {
            showMessage('Runner 删除成功！', 'success');
            loadRunners();
        } else {
            showMessage(`删除失败: ${data.message || response.statusText}`, 'error');
        }
    } catch (error) {
        showMessage(`删除失败: ${error.message}`, 'error');
    }
}

// Get status badge HTML
function getStatusBadge(status) {
    const statusMap = {
        'running': { class: 'status-running', text: '运行中' },
        'stopped': { class: 'status-stopped', text: '已停止' },
        'unknown': { class: 'status-unknown', text: '未知' }
    };

    const statusInfo = statusMap[status] || statusMap['unknown'];
    return `<span class="status-badge ${statusInfo.class}">${statusInfo.text}</span>`;
}

// Show message
function showMessage(text, type) {
    const messageDiv = document.getElementById('message');
    messageDiv.textContent = text;
    messageDiv.className = `message ${type}`;
    messageDiv.style.display = 'block';

    setTimeout(() => {
        messageDiv.style.display = 'none';
    }, 5000);
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return text.replace(/[&<>"']/g, m => map[m]);
}
