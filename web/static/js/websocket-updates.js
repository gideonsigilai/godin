/**
 * WebSocket Real-time Updates for Godin Framework
 * Handles automatic DOM updates when state changes occur
 */

class GodinWebSocketManager {
    constructor() {
        this.socket = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 1000;
        this.listeners = new Map();
        this.isConnected = false;
        this.pendingMessages = [];
        
        this.init();
    }

    init() {
        this.connect();
        this.setupEventListeners();
    }

    connect() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;
        
        try {
            this.socket = new WebSocket(wsUrl);
            this.setupSocketEventHandlers();
        } catch (error) {
            console.error('Failed to create WebSocket connection:', error);
            this.scheduleReconnect();
        }
    }

    setupSocketEventHandlers() {
        this.socket.onopen = (event) => {
            console.log('WebSocket connected');
            this.isConnected = true;
            this.reconnectAttempts = 0;
            
            // Send any pending messages
            while (this.pendingMessages.length > 0) {
                const message = this.pendingMessages.shift();
                this.socket.send(JSON.stringify(message));
            }
            
            // Subscribe to channels for elements that need updates
            this.subscribeToActiveChannels();
        };

        this.socket.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                this.handleMessage(data);
            } catch (error) {
                console.error('Failed to parse WebSocket message:', error);
            }
        };

        this.socket.onclose = (event) => {
            console.log('WebSocket disconnected');
            this.isConnected = false;
            
            if (!event.wasClean) {
                this.scheduleReconnect();
            }
        };

        this.socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    scheduleReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
            
            console.log(`Attempting to reconnect in ${delay}ms (attempt ${this.reconnectAttempts})`);
            
            setTimeout(() => {
                this.connect();
            }, delay);
        } else {
            console.error('Max reconnection attempts reached. Falling back to polling.');
            this.fallbackToPolling();
        }
    }

    handleMessage(data) {
        switch (data.type) {
            case 'value_change':
                this.handleValueChange(data);
                break;
            case 'state_update':
                this.handleStateUpdate(data);
                break;
            case 'widget_update':
                this.handleWidgetUpdate(data);
                break;
            default:
                console.warn('Unknown message type:', data.type);
        }
    }

    handleValueChange(data) {
        const { id, value, timestamp } = data;
        
        // Find all elements listening to this value notifier
        const elements = document.querySelectorAll(`[data-value-notifier-id="${id}"]`);
        
        elements.forEach(element => {
            this.updateElement(element, value, timestamp);
        });
        
        // Trigger custom event for other listeners
        const event = new CustomEvent('valueChanged', {
            detail: { id, value, timestamp }
        });
        document.body.dispatchEvent(event);
    }

    handleStateUpdate(data) {
        const { key, value } = data;
        
        // Find all elements listening to this state key
        const elements = document.querySelectorAll(`[data-state-key="${key}"]`);
        
        elements.forEach(element => {
            this.updateElement(element, value);
        });
    }

    handleWidgetUpdate(data) {
        const { widgetId, html } = data;
        
        const element = document.getElementById(widgetId);
        if (element) {
            element.innerHTML = html;
            
            // Re-initialize any interactive elements
            this.reinitializeInteractiveElements(element);
        }
    }

    updateElement(element, value, timestamp) {
        const updateMode = element.getAttribute('data-update-mode') || 'htmx';
        const debounce = parseInt(element.getAttribute('data-debounce')) || 0;
        
        // Apply debouncing if specified
        if (debounce > 0) {
            const lastUpdate = element.getAttribute('data-last-update');
            const now = Date.now();
            
            if (lastUpdate && (now - parseInt(lastUpdate)) < debounce) {
                return; // Skip update due to debouncing
            }
            
            element.setAttribute('data-last-update', now.toString());
        }
        
        switch (updateMode) {
            case 'websocket':
                this.updateElementDirectly(element, value);
                break;
            case 'htmx':
                this.triggerHTMXUpdate(element);
                break;
            case 'polling':
                // Polling is handled separately
                break;
        }
    }

    updateElementDirectly(element, value) {
        // Check if element has a builder function or should be updated directly
        const builderFunction = element.getAttribute('data-builder-function');
        
        if (builderFunction) {
            // Call the builder function to generate new content
            try {
                const newContent = window[builderFunction](value);
                element.innerHTML = newContent;
            } catch (error) {
                console.error('Failed to call builder function:', error);
            }
        } else {
            // Simple text update
            if (typeof value === 'string' || typeof value === 'number') {
                element.textContent = value;
            } else {
                element.textContent = JSON.stringify(value);
            }
        }
        
        // Re-initialize any interactive elements
        this.reinitializeInteractiveElements(element);
    }

    triggerHTMXUpdate(element) {
        // Trigger HTMX update if HTMX is available
        if (window.htmx) {
            htmx.trigger(element, 'valueChanged');
        }
    }

    reinitializeInteractiveElements(container) {
        // Re-initialize HTMX elements
        if (window.htmx) {
            htmx.process(container);
        }
        
        // Re-initialize any other interactive components
        const interactiveElements = container.querySelectorAll('[data-interactive]');
        interactiveElements.forEach(element => {
            this.initializeInteractiveElement(element);
        });
    }

    initializeInteractiveElement(element) {
        // Initialize interactive behaviors for the element
        const widgetType = element.getAttribute('data-widget-type');
        
        switch (widgetType) {
            case 'Button':
            case 'ElevatedButton':
            case 'TextButton':
            case 'OutlinedButton':
            case 'FilledButton':
            case 'IconButton':
            case 'FloatingActionButton':
                this.initializeButton(element);
                break;
            case 'Switch':
                this.initializeSwitch(element);
                break;
            case 'TextField':
            case 'TextFormField':
                this.initializeTextField(element);
                break;
        }
    }

    initializeButton(element) {
        // Add click handlers and visual feedback
        element.addEventListener('click', (e) => {
            this.addRippleEffect(element, e);
        });
    }

    initializeSwitch(element) {
        // Add change handlers for switch elements
        const input = element.querySelector('input[type=\"checkbox\"]');
        if (input) {
            input.addEventListener('change', (e) => {
                this.updateSwitchVisualState(element, e.target.checked);
            });
        }
    }

    initializeTextField(element) {
        // Add input handlers for text fields
        const input = element.querySelector('input, textarea');
        if (input) {
            input.addEventListener('input', (e) => {
                this.handleTextFieldInput(element, e.target.value);
            });
        }
    }

    addRippleEffect(element, event) {
        // Add Material Design ripple effect
        const ripple = document.createElement('span');
        const rect = element.getBoundingClientRect();
        const size = Math.max(rect.width, rect.height);
        const x = event.clientX - rect.left - size / 2;
        const y = event.clientY - rect.top - size / 2;
        
        ripple.style.width = ripple.style.height = size + 'px';
        ripple.style.left = x + 'px';
        ripple.style.top = y + 'px';
        ripple.classList.add('ripple');
        
        element.appendChild(ripple);
        
        setTimeout(() => {
            ripple.remove();
        }, 600);
    }

    updateSwitchVisualState(element, checked) {
        const thumb = element.querySelector('.godin-switch-thumb');
        if (thumb) {
            if (checked) {
                thumb.style.left = '22px';
            } else {
                thumb.style.left = '2px';
            }
        }
    }

    handleTextFieldInput(element, value) {
        // Handle text field input with debouncing
        const debounce = parseInt(element.getAttribute('data-debounce')) || 300;
        
        clearTimeout(element.inputTimeout);
        element.inputTimeout = setTimeout(() => {
            // Trigger any registered input handlers
            const event = new CustomEvent('textFieldInput', {
                detail: { element, value }
            });
            document.body.dispatchEvent(event);
        }, debounce);
    }

    subscribeToActiveChannels() {
        // Subscribe to channels for all elements that need WebSocket updates
        const elements = document.querySelectorAll('[data-ws-listen=\"true\"]');
        
        elements.forEach(element => {
            const channel = element.getAttribute('data-ws-channel');
            if (channel) {
                this.subscribe(channel);
            }
        });
    }

    subscribe(channel) {
        if (this.isConnected) {
            const message = {
                type: 'subscribe',
                channel: channel
            };
            this.socket.send(JSON.stringify(message));
        } else {
            this.pendingMessages.push({
                type: 'subscribe',
                channel: channel
            });
        }
    }

    unsubscribe(channel) {
        if (this.isConnected) {
            const message = {
                type: 'unsubscribe',
                channel: channel
            };
            this.socket.send(JSON.stringify(message));
        }
    }

    fallbackToPolling() {
        // Implement polling fallback for when WebSocket is unavailable
        const elements = document.querySelectorAll('[data-polling-endpoint]');
        
        elements.forEach(element => {
            const endpoint = element.getAttribute('data-polling-endpoint');
            const interval = parseInt(element.getAttribute('data-polling-interval')) || 5000;
            
            setInterval(() => {
                fetch(endpoint)
                    .then(response => response.json())
                    .then(data => {
                        this.updateElement(element, data.value);
                    })
                    .catch(error => {
                        console.error('Polling failed:', error);
                    });
            }, interval);
        });
    }

    setupEventListeners() {
        // Listen for DOM changes to initialize new elements
        const observer = new MutationObserver((mutations) => {
            mutations.forEach((mutation) => {
                mutation.addedNodes.forEach((node) => {
                    if (node.nodeType === Node.ELEMENT_NODE) {
                        this.initializeNewElements(node);
                    }
                });
            });
        });

        observer.observe(document.body, {
            childList: true,
            subtree: true
        });

        // Handle page visibility changes
        document.addEventListener('visibilitychange', () => {
            if (document.hidden) {
                // Page is hidden, reduce update frequency
                this.pauseUpdates();
            } else {
                // Page is visible, resume normal updates
                this.resumeUpdates();
            }
        });
    }

    initializeNewElements(container) {
        // Initialize WebSocket listeners for new elements
        const wsElements = container.querySelectorAll('[data-ws-listen=\"true\"]');
        wsElements.forEach(element => {
            const channel = element.getAttribute('data-ws-channel');
            if (channel) {
                this.subscribe(channel);
            }
        });

        // Initialize interactive elements
        const interactiveElements = container.querySelectorAll('[data-interactive]');
        interactiveElements.forEach(element => {
            this.initializeInteractiveElement(element);
        });
    }

    pauseUpdates() {
        // Reduce update frequency when page is not visible
        // This could involve unsubscribing from non-critical channels
    }

    resumeUpdates() {
        // Resume normal update frequency when page becomes visible
        this.subscribeToActiveChannels();
    }

    disconnect() {
        if (this.socket) {
            this.socket.close();
        }
    }
}

// Initialize WebSocket manager when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    window.godinWebSocket = new GodinWebSocketManager();
});

// Add CSS for ripple effects
const style = document.createElement('style');
style.textContent = `
    .ripple {
        position: absolute;
        border-radius: 50%;
        background-color: rgba(255, 255, 255, 0.6);
        transform: scale(0);
        animation: ripple-animation 0.6s linear;
        pointer-events: none;
    }

    @keyframes ripple-animation {
        to {
            transform: scale(4);
            opacity: 0;
        }
    }

    .godin-button, .godin-elevated-button, .godin-text-button, 
    .godin-outlined-button, .godin-filled-button, .godin-icon-button,
    .godin-floating-action-button {
        position: relative;
        overflow: hidden;
    }
`;
document.head.appendChild(style);