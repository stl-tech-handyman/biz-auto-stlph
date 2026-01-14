/**
 * Universal DataTables Auto-Initialization Script
 * Automatically applies DataTables.js to all tables on the page
 * with flexible configuration via data attributes
 */

(function() {
    'use strict';

    // Wait for DOM and jQuery/DataTables to be ready
    function initDataTables() {
        // Check if jQuery and DataTables are available
        if (typeof jQuery === 'undefined' || typeof jQuery.fn.DataTable === 'undefined') {
            console.warn('[DataTables Auto-Init] jQuery or DataTables not loaded. Skipping initialization.');
            return;
        }

        // Default configuration
        const defaultConfig = {
            pageLength: 25,
            lengthMenu: [[10, 25, 50, 100, -1], [10, 25, 50, 100, "All"]],
            order: [[0, 'asc']], // Sort by first column ascending
            dom: 'Bfrtip', // Buttons, filter, table, info, pagination
            buttons: ['copy', 'csv', 'excel', 'pdf', 'print'],
            language: {
                search: "Filter:",
                lengthMenu: "Show _MENU_ entries",
                info: "Showing _START_ to _END_ of _TOTAL_ entries",
                infoEmpty: "Showing 0 to 0 of 0 entries",
                infoFiltered: "(filtered from _TOTAL_ total entries)",
                zeroRecords: "No matching records found",
                paginate: {
                    first: "First",
                    last: "Last",
                    next: "Next",
                    previous: "Previous"
                }
            },
            responsive: true,
            autoWidth: false,
            scrollX: true
        };

        /**
         * Initialize DataTables on a table element
         */
        function initTable(table) {
            // Skip if already initialized
            if (jQuery.fn.DataTable.isDataTable(table)) {
                return;
            }

            // Skip if table has data-skip-datatables attribute
            if (table.hasAttribute('data-skip-datatables')) {
                return;
            }

            // Get custom configuration from data attributes
            const config = { ...defaultConfig };

            // Custom page length
            if (table.hasAttribute('data-page-length')) {
                config.pageLength = parseInt(table.getAttribute('data-page-length'), 10);
            }

            // Custom order (column index, direction)
            if (table.hasAttribute('data-order-column')) {
                const col = parseInt(table.getAttribute('data-order-column'), 10);
                const dir = table.getAttribute('data-order-dir') || 'asc';
                config.order = [[col, dir]];
            }

            // Disable buttons
            if (table.hasAttribute('data-no-buttons')) {
                config.buttons = [];
                config.dom = 'frtip'; // Remove 'B' (buttons)
            }

            // Custom buttons
            if (table.hasAttribute('data-buttons')) {
                config.buttons = table.getAttribute('data-buttons').split(',').map(b => b.trim());
            }

            // Disable pagination (show all)
            if (table.hasAttribute('data-no-pagination')) {
                config.paging = false;
            }

            // Disable search
            if (table.hasAttribute('data-no-search')) {
                config.searching = false;
                config.dom = config.dom.replace('f', ''); // Remove 'f' (filter)
            }

            // Enable column visibility toggle
            if (table.hasAttribute('data-column-toggle')) {
                if (!config.buttons.includes('colvis')) {
                    config.buttons.push('colvis');
                }
            }

            // Custom language/search placeholder
            if (table.hasAttribute('data-search-placeholder')) {
                config.language.search = table.getAttribute('data-search-placeholder');
            }

            // Grouping/Row grouping (requires DataTables RowGroup extension)
            if (table.hasAttribute('data-group-by')) {
                const groupCol = parseInt(table.getAttribute('data-group-by'), 10);
                config.rowGroup = {
                    dataSrc: groupCol
                };
            }

            // Initialize DataTable
            try {
                jQuery(table).DataTable(config);
                console.log('[DataTables Auto-Init] Initialized table:', table.id || table.className);
            } catch (error) {
                console.error('[DataTables Auto-Init] Error initializing table:', error, table);
            }
        }

        /**
         * Initialize all tables on the page
         */
        function initAllTables() {
            const tables = document.querySelectorAll('table.table:not(.no-datatables)');
            tables.forEach(table => {
                // Only init if table has tbody with content (not just loading state)
                const tbody = table.querySelector('tbody');
                if (tbody && tbody.children.length > 0) {
                    // Check if it's not just a loading row
                    const firstRow = tbody.children[0];
                    const isNotLoading = !firstRow.textContent.includes('Loading') && 
                                       !firstRow.textContent.includes('Running') &&
                                       !firstRow.classList.contains('loading-row');
                    
                    if (isNotLoading || tbody.children.length > 1) {
                        initTable(table);
                    }
                } else if (!tbody) {
                    // Table without tbody, init anyway
                    initTable(table);
                }
            });
        }

        // Initialize on page load
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', function() {
                setTimeout(initAllTables, 100); // Small delay for dynamic content
            });
        } else {
            setTimeout(initAllTables, 100);
        }

        // Watch for dynamically added tables
        const observer = new MutationObserver(function(mutations) {
            mutations.forEach(function(mutation) {
                mutation.addedNodes.forEach(function(node) {
                    if (node.nodeType === 1) { // Element node
                        // Check if the added node is a table
                        if (node.tagName === 'TABLE' && node.classList.contains('table')) {
                            setTimeout(() => initTable(node), 100);
                        }
                        // Check if the added node contains tables
                        const tables = node.querySelectorAll && node.querySelectorAll('table.table');
                        if (tables) {
                            tables.forEach(table => {
                                setTimeout(() => initTable(table), 100);
                            });
                        }
                    }
                });
            });
        });

        // Start observing
        observer.observe(document.body, {
            childList: true,
            subtree: true
        });

        // Expose reinit function for manual re-initialization
        window.reinitDataTables = function() {
            initAllTables();
        };
    }

    // Try to initialize immediately, or wait for dependencies
    if (typeof jQuery !== 'undefined' && typeof jQuery.fn.DataTable !== 'undefined') {
        initDataTables();
    } else {
        // Wait for dependencies to load
        let checkCount = 0;
        const checkInterval = setInterval(function() {
            checkCount++;
            if (typeof jQuery !== 'undefined' && typeof jQuery.fn.DataTable !== 'undefined') {
                clearInterval(checkInterval);
                initDataTables();
            } else if (checkCount > 50) { // 5 seconds max wait
                clearInterval(checkInterval);
                console.warn('[DataTables Auto-Init] Timeout waiting for dependencies');
            }
        }, 100);
    }
})();
