// Server Uptime Timer - Reusable component with Hot Reload
(function() {
    // Create timer HTML if it doesn't exist
    if (!document.getElementById('server-timer')) {
        const timerHTML = `
            <div id="server-timer">
                <div class="timer-time">
                    <div class="time-unit" id="timer-days" style="display: none;"><span class="prefix">D -</span> <span class="value">00</span></div>
                    <div class="time-unit" id="timer-hours"><span class="prefix">H -</span> <span class="value">00</span></div>
                    <div class="time-unit" id="timer-minutes"><span class="prefix">M -</span> <span class="value">00</span></div>
                    <div class="time-unit" id="timer-seconds"><span class="prefix">S -</span> <span class="value">00</span></div>
                </div>
            </div>
        `;
        document.body.insertAdjacentHTML('afterbegin', timerHTML);
    }

    const timerEl = document.getElementById('server-timer');
    const daysEl = document.getElementById('timer-days');
    const hoursEl = document.getElementById('timer-hours');
    const minutesEl = document.getElementById('timer-minutes');
    const secondsEl = document.getElementById('timer-seconds');
    
    let serverStartTime = null;
    let serverStartTimestamp = null; // Store the server start timestamp for hot reload detection
    let reloadCheckInterval = null;
    
    // Hot reload configuration
    const HOT_RELOAD_ENABLED = true; // Set to false to disable hot reload
    const HOT_RELOAD_CHECK_INTERVAL = 2000; // Check every 2 seconds
    
    // Periodic auto-reload configuration
    const AUTO_RELOAD_ENABLED = true; // Set to false to disable periodic auto-reload
    const AUTO_RELOAD_INTERVAL = 60000; // Reload every 60 seconds (1 minute)
    
    function updateDisplay(initialUptime) {
        let seconds;
        
        if (initialUptime !== undefined) {
            seconds = initialUptime;
        } else if (serverStartTime) {
            const now = Date.now();
            seconds = (now - serverStartTime) / 1000;
        } else {
            return;
        }
        
        const days = Math.floor(seconds / 86400);
        const hours = Math.floor((seconds % 86400) / 3600);
        const minutes = Math.floor((seconds % 3600) / 60);
        const secs = Math.floor(seconds % 60);
        
        if (days > 0) {
            daysEl.style.display = 'flex';
            daysEl.querySelector('.value').textContent = String(days).padStart(2, '0');
        } else {
            daysEl.style.display = 'none';
        }
        
        hoursEl.querySelector('.value').textContent = String(hours).padStart(2, '0');
        minutesEl.querySelector('.value').textContent = String(minutes).padStart(2, '0');
        secondsEl.querySelector('.value').textContent = String(secs).padStart(2, '0');
    }
    
    async function syncServerUptime() {
        try {
            const response = await fetch('/api/health');
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            const data = await response.json();
            const uptimeSeconds = data.uptime;
            if (!serverStartTime) {
                const serverTimestamp = new Date(data.timestamp);
                serverStartTime = serverTimestamp.getTime() - (uptimeSeconds * 1000);
            }
            timerEl.classList.remove('error');
            updateDisplay(uptimeSeconds);
        } catch (error) {
            console.error('Failed to sync server uptime:', error);
            timerEl.classList.add('error');
        }
    }
    
    // Hot reload: Check if server has restarted
    async function checkServerRestart() {
        if (!HOT_RELOAD_ENABLED) return;
        
        try {
            const response = await fetch('/api/server-start-time');
            if (!response.ok) {
                // Server might be restarting, wait and retry
                return;
            }
            
            const data = await response.json();
            const currentStartTimestamp = data.timestamp;
            
            // If we haven't stored a timestamp yet, store it now
            if (serverStartTimestamp === null) {
                serverStartTimestamp = currentStartTimestamp;
                console.log('[Hot Reload] Initial server start time:', new Date(currentStartTimestamp * 1000).toLocaleString());
                return;
            }
            
            // If the timestamp has changed, server has restarted - reload the page
            if (currentStartTimestamp !== serverStartTimestamp) {
                console.log('[Hot Reload] Server restarted detected! Reloading page...');
                console.log('[Hot Reload] Old start time:', new Date(serverStartTimestamp * 1000).toLocaleString());
                console.log('[Hot Reload] New start time:', new Date(currentStartTimestamp * 1000).toLocaleString());
                
                // Small delay to ensure server is fully ready
                setTimeout(() => {
                    window.location.reload();
                }, 500);
            }
        } catch (error) {
            // Silently handle errors (server might be restarting)
            // Don't log to avoid console spam
        }
    }
    
    // Initialize
    syncServerUptime();
    setInterval(() => updateDisplay(), 1000);
    setInterval(syncServerUptime, 30000);
    
    // Start hot reload checking
    if (HOT_RELOAD_ENABLED) {
        // Initial check after a short delay
        setTimeout(checkServerRestart, 1000);
        // Then check periodically
        reloadCheckInterval = setInterval(checkServerRestart, HOT_RELOAD_CHECK_INTERVAL);
        console.log('[Hot Reload] Enabled - checking every', HOT_RELOAD_CHECK_INTERVAL / 1000, 'seconds');
    }
    
    // Start periodic auto-reload
    if (AUTO_RELOAD_ENABLED) {
        console.log('[Auto Reload] Enabled - page will reload every', AUTO_RELOAD_INTERVAL / 1000, 'seconds');
        setInterval(() => {
            console.log('[Auto Reload] Reloading page...');
            window.location.reload();
        }, AUTO_RELOAD_INTERVAL);
    }
})();
