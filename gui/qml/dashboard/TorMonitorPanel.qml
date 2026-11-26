import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

Rectangle {
    id: root

    required property var torMonitor

    readonly property int headerH: 40
    readonly property int bodyH:   256
    property bool expanded: false

    implicitHeight: root.expanded ? root.headerH + root.bodyH : root.headerH
    clip: true
    color: Theme.bgPanel
    border.width: 1
    border.color: Theme.borderLight

    Behavior on implicitHeight {
        NumberAnimation { duration: Theme.durationMedium; easing.type: Easing.InOutQuad }
    }

    // ── Header ───────────────────────────────────────────────────────────────
    Rectangle {
        id: header
        anchors { top: parent.top; left: parent.left; right: parent.right }
        height: root.headerH
        color: "transparent"

        RowLayout {
            anchors { fill: parent; leftMargin: Theme.spacingMedium; rightMargin: Theme.spacingMedium }
            spacing: Theme.spacingSmall

            // Expand/collapse arrow
            Text {
                text: root.expanded ? "▼" : "▶"
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: 10
            }

            Text {
                text: "TOR MONITOR"
                color: Theme.textPrimary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
                font.letterSpacing: 1.2
            }

            // Status dot + text
            Rectangle {
                width: 7; height: 7; radius: 4
                color: root.torMonitor.connected ? Theme.success : Theme.textDisabled
            }

            Text {
                text: root.torMonitor.status
                color: root.torMonitor.connected ? Theme.success : Theme.textDisabled
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
            }

            Item { Layout.fillWidth: true }

            // Current rates
            Text {
                visible: root.torMonitor.connected
                text: "↓ " + formatRate(root.torMonitor.downRate)
                     + "   ↑ " + formatRate(root.torMonitor.upRate)
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize

                function formatRate(b) {
                    if (b >= 1048576) return (b / 1048576).toFixed(1) + " MB/s"
                    if (b >= 1024)    return (b / 1024).toFixed(1)    + " KB/s"
                    return b.toFixed(0) + " B/s"
                }
            }
        }

        MouseArea {
            anchors.fill: parent
            cursorShape: Qt.PointingHandCursor
            onClicked: root.expanded = !root.expanded
        }
    }

    // ── Body ─────────────────────────────────────────────────────────────────
    ColumnLayout {
        anchors {
            top: header.bottom
            left: parent.left; right: parent.right; bottom: parent.bottom
            margins: Theme.spacingMedium
            topMargin: Theme.spacingSmall
        }
        spacing: Theme.spacingSmall
        visible: root.expanded

        // Rate labels + legend
        RowLayout {
            Layout.fillWidth: true
            spacing: Theme.spacingMedium

            Rectangle { width: 20; height: 2; color: Theme.success }
            Text {
                text: "Download"
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
            }

            Rectangle { width: 20; height: 2; color: Theme.accent }
            Text {
                text: "Upload"
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
            }

            Item { Layout.fillWidth: true }

            Text {
                visible: root.torMonitor.circuitCount > 0
                text: root.torMonitor.circuitCount + " circuit" + (root.torMonitor.circuitCount !== 1 ? "s" : "")
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
            }
        }

        // Bandwidth graph
        BandwidthGraph {
            Layout.fillWidth: true
            height: 100
            downHistory: root.torMonitor.downHistory
            upHistory:   root.torMonitor.upHistory

            Rectangle {
                anchors.fill: parent
                color: "transparent"
                border.width: 1
                border.color: Theme.borderLight
            }
        }

        // No data placeholder
        Item {
            Layout.fillWidth: true
            height: 24
            visible: !root.torMonitor.connected

            Text {
                anchors.centerIn: parent
                text: {
                    const s = root.torMonitor.status
                    if (s === "Connecting...") return "Connecting to Tor control port..."
                    if (s.indexOf("auth") >= 0) return "Cannot authenticate — log out and log in after install to apply group membership"
                    if (s.indexOf("Cannot connect") >= 0) return "Tor control port unavailable — start Anonimax first"
                    return "Tor is not active — start Anonimax to see live data"
                }
                color: Theme.textDisabled
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
            }
        }

        // Circuit list
        ListView {
            id: circuitList
            Layout.fillWidth: true
            height: Math.min(contentHeight, 72)
            clip: true
            visible: root.torMonitor.circuitCount > 0
            model: root.torMonitor.circuits

            delegate: Text {
                required property var modelData
                width: circuitList.width
                text: "● " + (modelData.path || "")
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
                elide: Text.ElideRight
            }
        }
    }
}
