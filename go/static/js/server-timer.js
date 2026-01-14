// Server Uptime Timer - Reusable component with Hot Reload and Auto-Restart
(function() {
    // Create timer HTML if it doesn't exist
    if (!document.getElementById('server-timer')) {
        const timerHTML = `
            <div id="server-timer">
                <div class="timer-content">
                    <div class="timer-header">
                        <span class="timer-title">Server Uptime</span>
                    </div>
                    <div class="timer-time">
                        <div class="time-unit" id="timer-days" style="display: none;"><span class="prefix">D -</span> <span class="value">00</span></div>
                        <div class="time-unit" id="timer-hours"><span class="prefix">H -</span> <span class="value">00</span></div>
                        <div class="time-unit" id="timer-minutes"><span class="prefix">M -</span> <span class="value">00</span></div>
                        <div class="time-unit" id="timer-seconds"><span class="prefix">S -</span> <span class="value">00</span></div>
                    </div>
                    <div class="restart-timer-section">
                        <div class="restart-timer-header">
                            <span class="restart-timer-label">Auto Restart</span>
                            <button id="restart-timer-toggle" class="btn-toggle" title="Pause/Resume">
                                <i class="bi bi-pause-fill"></i>
                            </button>
                        </div>
                        <div class="restart-countdown" id="restart-countdown">
                            <span class="countdown-label">Next restart in:</span>
                            <span class="countdown-time" id="countdown-time">--:--</span>
                        </div>
                        <div class="restart-interval-selector">
                            <label>Interval:</label>
                            <select id="restart-interval-select">
                                <option value="1">1 min</option>
                                <option value="2">2 min</option>
                                <option value="3">3 min</option>
                                <option value="5" selected>5 min</option>
                                <option value="10">10 min</option>
                            </select>
                        </div>
                        <button id="manual-restart-btn" class="btn-manual-restart" title="Manually restart server">
                            <i class="bi bi-arrow-clockwise"></i> Restart Now
                        </button>
                    </div>
                </div>
                <button id="timer-collapse-btn" class="timer-collapse-btn" title="Collapse/Expand">
                    <i class="bi bi-chevron-right"></i>
                </button>
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
    
    // Auto-restart configuration
    let autoRestartEnabled = true; // Default enabled
    let autoRestartInterval = 5 * 60 * 1000; // Default 5 minutes in milliseconds
    let restartCountdownInterval = null;
    let nextRestartTime = null;
    let isPaused = false;
    
    // Collapse state
    let isCollapsed = false;
    
    // Load saved settings from localStorage
    const savedInterval = localStorage.getItem('autoRestartInterval');
    const savedPaused = localStorage.getItem('autoRestartPaused');
    const savedCollapsed = localStorage.getItem('timerCollapsed');
    if (savedInterval) {
        autoRestartInterval = parseInt(savedInterval) * 60 * 1000;
        document.getElementById('restart-interval-select').value = savedInterval;
    }
    if (savedPaused === 'true') {
        isPaused = true;
        autoRestartEnabled = false;
        updateToggleButton();
    }
    if (savedCollapsed === 'true') {
        isCollapsed = true;
        timerEl.classList.add('collapsed');
    }
    // Always update button state on initialization
    updateCollapseButton();
    
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
            // Use /api/time to get current server time, then calculate from server start time
            const [timeResponse, startTimeResponse] = await Promise.all([
                fetch('/api/time'),
                fetch('/api/server-start-time')
            ]);
            
            if (!timeResponse.ok || !startTimeResponse.ok) {
                throw new Error(`HTTP ${timeResponse.status} or ${startTimeResponse.status}`);
            }
            
            const timeData = await timeResponse.json();
            const startTimeData = await startTimeResponse.json();
            
            // Calculate uptime from server start timestamp
            const serverCurrentTime = timeData.unix * 1000; // Convert to milliseconds
            const serverStartUnix = startTimeData.timestamp * 1000; // Convert to milliseconds
            const uptimeSeconds = (serverCurrentTime - serverStartUnix) / 1000;
            
            if (!serverStartTime) {
                serverStartTime = serverStartUnix;
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
                // Initialize restart timer after first sync
                if (autoRestartEnabled && !isPaused) {
                    scheduleNextRestart();
                }
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
    
    // Auto-restart functionality
    function scheduleNextRestart() {
        if (isPaused || !autoRestartEnabled) return;
        
        nextRestartTime = Date.now() + autoRestartInterval;
        updateCountdown();
        
        if (restartCountdownInterval) {
            clearInterval(restartCountdownInterval);
        }
        
        restartCountdownInterval = setInterval(() => {
            if (isPaused || !autoRestartEnabled) {
                clearInterval(restartCountdownInterval);
                return;
            }
            
            const now = Date.now();
            if (now >= nextRestartTime) {
                triggerServerRestart();
            } else {
                updateCountdown();
            }
        }, 1000); // Update every second
    }
    
    function updateCountdown() {
        if (!nextRestartTime || isPaused || !autoRestartEnabled) {
            document.getElementById('countdown-time').textContent = '--:--';
            return;
        }
        
        const now = Date.now();
        const remaining = Math.max(0, nextRestartTime - now);
        const minutes = Math.floor(remaining / 60000);
        const seconds = Math.floor((remaining % 60000) / 1000);
        
        document.getElementById('countdown-time').textContent = 
            `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
    }
    
    async function triggerServerRestart() {
        console.log('[Auto Restart] Triggering server restart...');
        
        try {
            const response = await fetch('/api/server/restart', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            
            if (response.ok) {
                const data = await response.json();
                console.log('[Auto Restart] Server restart initiated:', data);
                // Wait a bit then reload page
                setTimeout(() => {
                    window.location.reload();
                }, 2000);
            } else {
                console.error('[Auto Restart] Failed to restart server:', response.status);
            }
        } catch (error) {
            console.error('[Auto Restart] Error triggering restart:', error);
        }
    }
    
    async function manualRestart() {
        if (!confirm('Are you sure you want to manually restart the server?')) {
            return;
        }
        
        console.log('[Manual Restart] Triggering server restart...');
        // #region agent log
        fetch('http://127.0.0.1:7242/ingest/884502d0-ec5e-4430-9e0a-a08e454d487b',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'server-timer.js:250',message:'manualRestart called',data:{},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'A'})}).catch(()=>{});
        // #endregion
        
        const btn = document.getElementById('manual-restart-btn');
        const originalHTML = btn.innerHTML;
        btn.disabled = true;
        btn.innerHTML = '<i class="bi bi-hourglass-split"></i> Restarting...';
        
        try {
            // #region agent log
            fetch('http://127.0.0.1:7242/ingest/884502d0-ec5e-4430-9e0a-a08e454d487b',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'server-timer.js:262',message:'fetch request starting',data:{url:'/api/server/restart',method:'POST'},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'B'})}).catch(()=>{});
            // #endregion
            
            const response = await fetch('/api/server/restart', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            
            // #region agent log
            fetch('http://127.0.0.1:7242/ingest/884502d0-ec5e-4430-9e0a-a08e454d487b',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'server-timer.js:272',message:'response received',data:{ok:response.ok,status:response.status,statusText:response.statusText,bodyUsed:response.bodyUsed},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'C'})}).catch(()=>{});
            // #endregion
            
            // Clone response to avoid "body stream already read" error
            const responseClone = response.clone();
            
            if (response.ok) {
                // #region agent log
                fetch('http://127.0.0.1:7242/ingest/884502d0-ec5e-4430-9e0a-a08e454d487b',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'server-timer.js:278',message:'response.ok is true, reading json',data:{},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'A'})}).catch(()=>{});
                // #endregion
                const data = await response.json();
                // #region agent log
                fetch('http://127.0.0.1:7242/ingest/884502d0-ec5e-4430-9e0a-a08e454d487b',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'server-timer.js:281',message:'json parsed successfully',data:{status:data.status,message:data.message},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'A'})}).catch(()=>{});
                // #endregion
                console.log('[Manual Restart] Server restart initiated:', data);
                // Wait a bit then reload page
                setTimeout(() => {
                    window.location.reload();
                }, 2000);
            } else {
                // #region agent log
                fetch('http://127.0.0.1:7242/ingest/884502d0-ec5e-4430-9e0a-a08e454d487b',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'server-timer.js:290',message:'response.ok is false, reading error body',data:{status:response.status,bodyUsed:response.bodyUsed},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'B'})}).catch(()=>{});
                // #endregion
                // Try to get error message from response - use clone to avoid stream consumption issues
                let errorMessage = `HTTP ${response.status}: ${response.statusText}`;
                try {
                    const errorData = await responseClone.json();
                    // #region agent log
                    fetch('http://127.0.0.1:7242/ingest/884502d0-ec5e-4430-9e0a-a08e454d487b',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'server-timer.js:295',message:'error json parsed',data:{error:errorData.error,message:errorData.message},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'B'})}).catch(()=>{});
                    // #endregion
                    if (errorData.error || errorData.message) {
                        errorMessage = errorData.error || errorData.message;
                    }
                } catch (e) {
                    // #region agent log
                    fetch('http://127.0.0.1:7242/ingest/884502d0-ec5e-4430-9e0a-a08e454d487b',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'server-timer.js:301',message:'error reading json, trying text',data:{error:e.message,errorName:e.name},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'C'})}).catch(()=>{});
                    // #endregion
                    // If response is not JSON, try to read as text using the clone (not original response)
                    try {
                        const text = await responseClone.text();
                        // #region agent log
                        fetch('http://127.0.0.1:7242/ingest/884502d0-ec5e-4430-9e0a-a08e454d487b',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'server-timer.js:305',message:'text read successfully',data:{textLength:text.length},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'C'})}).catch(()=>{});
                        // #endregion
                        if (text) {
                            errorMessage = text;
                        }
                    } catch (textError) {
                        // #region agent log
                        fetch('http://127.0.0.1:7242/ingest/884502d0-ec5e-4430-9e0a-a08e454d487b',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'server-timer.js:310',message:'failed to read text',data:{error:textError.message,errorName:textError.name},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'C'})}).catch(()=>{});
                        // #endregion
                        // If both fail, just use the status text we already have
                        console.warn('[Manual Restart] Could not read error body:', textError);
                    }
                }
                
                console.error('[Manual Restart] Failed to restart server:', errorMessage);
                btn.disabled = false;
                btn.innerHTML = originalHTML;
                alert('Failed to restart server:\n\n' + errorMessage + '\n\nNote: Server restart is only available in development environment (ENV=dev or ENV not set).');
            }
        } catch (error) {
            // #region agent log
            fetch('http://127.0.0.1:7242/ingest/884502d0-ec5e-4430-9e0a-a08e454d487b',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({location:'server-timer.js:322',message:'catch block executed',data:{error:error.message,errorName:error.name,stack:error.stack},timestamp:Date.now(),sessionId:'debug-session',runId:'run1',hypothesisId:'A'})}).catch(()=>{});
            // #endregion
            console.error('[Manual Restart] Error triggering restart:', error);
            btn.disabled = false;
            btn.innerHTML = originalHTML;
            alert('Error restarting server: ' + error.message + '\n\nCheck the browser console for more details.');
        }
    }
    
    function toggleCollapse() {
        try {
            isCollapsed = !isCollapsed;
            timerEl.classList.toggle('collapsed', isCollapsed);
            localStorage.setItem('timerCollapsed', isCollapsed.toString());
            updateCollapseButton();
        } catch (error) {
            console.error('Error toggling collapse:', error);
        }
    }
    
    function updateCollapseButton() {
        try {
            const btn = document.getElementById('timer-collapse-btn');
            if (!btn) {
                console.warn('Collapse button not found');
                return;
            }
            const icon = btn.querySelector('i');
            if (!icon) {
                console.warn('Collapse button icon not found');
                return;
            }
            if (isCollapsed) {
                icon.className = 'bi bi-chevron-left';
                btn.title = 'Expand';
            } else {
                icon.className = 'bi bi-chevron-right';
                btn.title = 'Collapse';
            }
        } catch (error) {
            console.error('Error updating collapse button:', error);
        }
    }
    
    function updateToggleButton() {
        const toggleBtn = document.getElementById('restart-timer-toggle');
        const icon = toggleBtn.querySelector('i');
        if (isPaused) {
            icon.className = 'bi bi-play-fill';
            toggleBtn.title = 'Resume';
            toggleBtn.classList.add('paused');
        } else {
            icon.className = 'bi bi-pause-fill';
            toggleBtn.title = 'Pause';
            toggleBtn.classList.remove('paused');
        }
    }
    
    // Event listeners
    document.getElementById('restart-timer-toggle').addEventListener('click', () => {
        isPaused = !isPaused;
        autoRestartEnabled = !isPaused;
        localStorage.setItem('autoRestartPaused', isPaused.toString());
        updateToggleButton();
        
        if (isPaused) {
            if (restartCountdownInterval) {
                clearInterval(restartCountdownInterval);
            }
            document.getElementById('countdown-time').textContent = '--:--';
        } else {
            scheduleNextRestart();
        }
    });
    
    document.getElementById('restart-interval-select').addEventListener('change', (e) => {
        const minutes = parseInt(e.target.value);
        autoRestartInterval = minutes * 60 * 1000;
        localStorage.setItem('autoRestartInterval', minutes.toString());
        
        if (!isPaused && autoRestartEnabled) {
            scheduleNextRestart();
        }
    });
    
    // Collapse button event listener - stop propagation to prevent conflicts
    const collapseBtn = document.getElementById('timer-collapse-btn');
    if (collapseBtn) {
        collapseBtn.addEventListener('click', (e) => {
            e.stopPropagation(); // Prevent event from bubbling to timer element
            toggleCollapse();
        });
    }
    
    const manualRestartBtn = document.getElementById('manual-restart-btn');
    if (manualRestartBtn) {
        manualRestartBtn.addEventListener('click', manualRestart);
    }
    
    // Allow clicking anywhere on the collapsed timer to expand it (but not on buttons)
    timerEl.addEventListener('click', (e) => {
        // Don't trigger if clicking on a button or inside timer-content
        if (e.target.closest('button') || e.target.closest('.timer-content')) {
            return;
        }
        if (isCollapsed) {
            // Only expand if clicking on the collapsed timer itself, not on content or buttons
            toggleCollapse();
        }
    });
    
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
    
    // Start auto-restart if enabled
    if (autoRestartEnabled && !isPaused) {
        // Wait for server start time to be synced before scheduling
        setTimeout(() => {
            if (serverStartTimestamp !== null) {
                scheduleNextRestart();
            }
        }, 2000);
    }
})();
