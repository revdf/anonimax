import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import QtQuick.Window

ApplicationWindow {
    id: window

    required property var backend
    required property var desktop
    required property var config
    property bool safeKill: false
    property string _lastState: ""
    property string _lastProgress: ""

    property string logText: ""

    function logEntry(message, entryType) {
        var prefix = entryType === "error"   ? "✗"
                   : entryType === "success" ? "✓"
                   : entryType === "warn"    ? "⚠"
                   : "→"
        logText += Qt.formatTime(new Date(), "hh:mm:ss") + "  " + prefix + " " + message + "\n"
    }
    // The tray needs fresh status even while the window is hidden.
    readonly property bool live: window.visible
                                 && window.visibility !== Window.Minimized
    // The window opens exactly as tall as its content and cannot grow past it,
    // so no space is ever left blank; making it shorter is allowed and the
    // view scrolls instead.
    readonly property int contentMargins: 2 * (window.width < 600
                                               ? 12 : Theme.contentMargin)
    readonly property int naturalHeight: Math.ceil(contentColumn.implicitHeight
                                                   + headerBar.implicitHeight)
                                         + Theme.spacingMedium
                                         + window.contentMargins
    width: 840
    height: window.naturalHeight
    minimumWidth: 480
    minimumHeight: 240
    maximumHeight: window.naturalHeight
    visible: true
    title: "Anonimax"
    color: Theme.bgWindow
    font.family: Theme.uiFont

    Component.onCompleted: {
        window.logEntry("Anonimax GUI iniciada — a conectar ao daemon...", "info")
        window.backend.refresh()
    }
    onClosing: close => {
        aboutWindow.close();
        if (logPanel.saveOnExit && window.logText !== "")
            window.config.saveLog(window.logText)
        if (window.desktop.hideToTray())
            close.accepted = false;
    }

    onLiveChanged: if (window.live) window.backend.refresh()

    Timer {
        interval: 3000
        running: window.live || window.desktop.trayAvailable
        repeat: true
        onTriggered: window.backend.refresh()
    }

    Rectangle {
        anchors.fill: parent
        gradient: Gradient {
            GradientStop { position: 0; color: Theme.bgWindow }
            GradientStop { position: 0.62; color: "#071923" }
            GradientStop { position: 1; color: Theme.bgBottom }
        }
    }

    CrtOverlay { anchors.fill: parent }

    ColumnLayout {
        anchors.fill: parent
        anchors.margins: window.width < 600 ? 12 : Theme.contentMargin
        spacing: Theme.spacingMedium

        HeaderBar {
            id: headerBar

            Layout.fillWidth: true
            onAboutRequested: aboutWindow.open()
        }

        ScrollView {
            id: contentScroll

            Layout.fillWidth: true
            Layout.fillHeight: true
            contentWidth: availableWidth
            clip: true

            ScrollBar.horizontal.policy: ScrollBar.AlwaysOff
            ScrollBar.vertical.policy: ScrollBar.AsNeeded

            Column {
                id: contentColumn

                width: contentScroll.availableWidth
                spacing: Theme.spacingMedium

                ProtectionPanel {
                    width: parent.width
                    backend: window.backend
                    safeKill: window.safeKill
                    onSafeKillToggled: enabled => {
                        if (enabled)
                            safeKillDialog.open();
                        else
                            window.safeKill = false;
                    }
                }

                ConfigPanel {
                    width: parent.width
                    config: window.config
                    onProxyNotConfigured: proxyWarningDialog.open()
                    onVpnNotConfigured: vpnWarningDialog.open()
                }

                ConnectionOverview {
                    width: parent.width
                    backend: window.backend
                }

                ActionPanel {
                    id: actionPanel
                    width: parent.width
                    backend: window.backend
                    safeKill: window.safeKill
                    onPanicRequested: panicDialog.open()
                }

                LogPanel {
                    id: logPanel
                    width: parent.width
                    logText: window.logText
                    onClearRequested: window.logText = ""
                }
            }
        }
    }

    AboutWindow {
        id: aboutWindow
        transientParent: window
    }

    ConfirmDialog {
        id: panicDialog

        badge: "PANIC"
        confirmText: "Execute Panic"
        message: "This will immediately cut all network connections, stop Tor, "
                 + "clear traces and POWER OFF the machine. All unsaved work will be lost. "
                 + "This cannot be undone."
        onAccepted: window.backend.runPanic()
    }

    ConfirmDialog {
        id: safeKillDialog

        badge: "SAFEKILL"
        confirmText: "Enable SafeKill"
        message: "Enabling SafeKill does not run it now. Before each start, stop "
                 + "or restart, it requests termination of selected applications, "
                 + "then clears selected caches, history and sessions with BleachBit. "
                 + "Unsaved work may be lost and cleared data cannot be recovered."
        onAccepted: window.safeKill = true
    }

    ErrorDialog {
        id: errorDialog

        backend: window.backend
    }

    ConfirmDialog {
        id: proxyWarningDialog

        badge: "PROXY"
        accentColor: Theme.accent
        confirmText: "Configure"
        message: "No proxy is configured. You need a SOCKS5 or HTTP proxy "
                 + "already running on your machine before enabling this option.\n\n"
                 + "Do you want to configure one now?"
        onAccepted: proxyConfigDialog.open()
    }

    ProxyConfigDialog {
        id: proxyConfigDialog

        config: window.config
    }

    ConfirmDialog {
        id: vpnWarningDialog

        badge: "VPN"
        accentColor: Theme.accent
        confirmText: "Configure"
        message: "No VPN configuration file is set. You need an OpenVPN (.ovpn) "
                 + "or WireGuard (.conf) config file on this machine.\n\n"
                 + "Do you want to set the path now?"
        onAccepted: vpnConfigDialog.open()
    }

    VpnConfigDialog {
        id: vpnConfigDialog

        config: window.config
    }

    Connections {
        target: window.desktop

        function onLifecycleRequested(action: string) {
            if (action === "start")
                window.backend.start(window.safeKill);
            else if (action === "stop")
                window.backend.stop(window.safeKill);
            else if (action === "restart")
                window.backend.restart(window.safeKill);
        }
    }

    Connections {
        target: window.backend

        function onChanged() {
            if (window.backend.errorMessage !== "" && !errorDialog.opened) {
                window.logEntry("✗ Error: " + window.backend.errorMessage, "error")
                errorDialog.open();
            }

            var state = window.backend.protectionState
            if (state !== window._lastState) {
                window._lastState = state
                var t = state === "ROUTED"     ? "success"
                      : state === "INACTIVE"   ? "info"
                      : state === "UNAVAILABLE" ? "warn"
                      : "warn"
                window.logEntry("Estado → " + state, t)
            }

            if (window.backend.busy) {
                var prog = window.backend.progress
                if (prog !== window._lastProgress) {
                    window._lastProgress = prog
                    window.logEntry(prog, "info")
                }
            }
        }

        function onOperationFinished(action, success, message) {
            window.logEntry((success ? "✓ " : "✗ ") + action + ": " + message,
                            success ? "success" : "error")
        }
    }
}
