import QtQuick
import QtQuick.Layouts
import QtQuick.Controls

Rectangle {
    id: root

    required property string logText
    property bool saveOnExit: false
    property bool expanded: false

    signal clearRequested()

    readonly property int headerH: 36
    readonly property int logH: 200

    implicitHeight: root.expanded ? root.headerH + root.logH : root.headerH
    clip: true
    color: Theme.bgPanel
    border.width: 1
    border.color: Theme.borderLight

    Behavior on implicitHeight {
        NumberAnimation { duration: Theme.durationMedium; easing.type: Easing.InOutQuad }
    }

    // ── Header (always visible) ──────────────────────────────────────
    RowLayout {
        id: header

        anchors.left: parent.left
        anchors.right: parent.right
        anchors.top: parent.top
        height: root.headerH
        anchors.leftMargin: Theme.spacingMedium
        anchors.rightMargin: Theme.spacingMedium
        spacing: Theme.spacingSmall

        // Arrow + label — click to toggle
        RowLayout {
            spacing: Theme.spacingSmall

            MouseArea {
                anchors.fill: parent
                cursorShape: Qt.PointingHandCursor
                onClicked: root.expanded = !root.expanded
            }

            Text {
                text: root.expanded ? "▼" : "▶"
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: Theme.badgeSize
            }

            Text {
                text: "ACTIVITY LOG"
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: Theme.badgeSize
                font.weight: Font.Medium
            }
        }

        Item { Layout.fillWidth: true }

        // Save on exit checkbox
        RowLayout {
            spacing: Theme.spacingTiny

            CheckBox {
                id: saveCheck
                checked: root.saveOnExit
                onCheckedChanged: root.saveOnExit = checked
                padding: 0

                indicator: Rectangle {
                    width: 14
                    height: 14
                    color: saveCheck.checked ? Theme.accent : Theme.buttonSecondary
                    border.width: 1
                    border.color: saveCheck.checked ? Theme.accent : Theme.buttonBorder

                    Text {
                        anchors.centerIn: parent
                        text: "✓"
                        color: Theme.textOnAccent
                        font.pixelSize: 10
                        visible: saveCheck.checked
                    }
                }
            }

            Text {
                text: "Guardar ao fechar"
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: Theme.badgeSize
            }
        }

        // Clear — only when expanded
        Text {
            text: "LIMPAR"
            color: Theme.accent
            font.family: Theme.uiFont
            font.pixelSize: Theme.badgeSize
            font.bold: true
            visible: root.expanded

            MouseArea {
                anchors.fill: parent
                cursorShape: Qt.PointingHandCursor
                onClicked: root.clearRequested()
            }
        }
    }

    // ── Separator ────────────────────────────────────────────────────
    Rectangle {
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.top: header.bottom
        height: 1
        color: Theme.borderLight
        visible: root.expanded
    }

    // ── Log area (visible when expanded) ─────────────────────────────
    ScrollView {
        id: scroll

        anchors.left: parent.left
        anchors.right: parent.right
        anchors.top: header.bottom
        anchors.topMargin: 1
        anchors.bottom: parent.bottom
        anchors.margins: Theme.spacingSmall
        visible: root.expanded
        clip: true
        ScrollBar.horizontal.policy: ScrollBar.AlwaysOff
        ScrollBar.vertical.policy: ScrollBar.AsNeeded

        TextArea {
            readOnly: true
            text: root.logText
            color: Theme.textPrimary
            font.family: "monospace"
            font.pixelSize: Theme.detailSize
            background: null
            wrapMode: TextArea.Wrap
            padding: 0
            selectByMouse: true

            onTextChanged: Qt.callLater(() => scroll.ScrollBar.vertical.position = 1.0)
        }
    }
}
