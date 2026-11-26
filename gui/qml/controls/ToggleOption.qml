import QtQuick
import QtQuick.Controls

AbstractButton {
    id: root

    property alias label: textLabel.text
    property string tooltip: ""
    property color accentColor: Theme.accent

    checkable: false
    hoverEnabled: true
    opacity: enabled ? 1 : 0.5
    implicitHeight: Theme.buttonHeight
    implicitWidth: textLabel.implicitWidth + 14 + Theme.spacingSmall * 2 + Theme.spacingLarge * 2
    Accessible.name: label

    contentItem: Item {
        anchors.fill: parent
        transform: Translate {
            x: root.down ? Theme.buttonShadowSize : 0
            y: root.down ? Theme.buttonShadowSize : 0
        }

        Row {
            anchors.centerIn: parent
            anchors.left: parent.left
            anchors.right: parent.right
            anchors.leftMargin: Theme.spacingLarge
            spacing: Theme.spacingSmall

            Icon {
                icon: root.checked ? "check" : "close"
                color: root.checked ? Theme.textOnAccent : Theme.textSecondary
                anchors.verticalCenter: parent.verticalCenter
                width: 14
                height: 14
            }

            Text {
                id: textLabel

                color: root.checked
                       ? Theme.textOnAccent
                       : (root.hovered ? Theme.textPrimary : Theme.textSecondary)
                font.family: Theme.uiFont
                font.pixelSize: Theme.buttonSize
                font.bold: root.checked
                anchors.verticalCenter: parent.verticalCenter
                elide: Text.ElideRight
                Behavior on color {
                    ColorAnimation { duration: Theme.durationFast }
                }
            }
        }
    }

    background: Item {
        Rectangle {
            x: root.down ? 0 : Theme.buttonShadowSize
            y: root.down ? 0 : Theme.buttonShadowSize
            width: parent.width
            height: parent.height
            color: root.checked ? Theme.buttonShadow : Theme.buttonSecondaryShadow
        }

        Rectangle {
            id: bg
            x: root.down ? Theme.buttonShadowSize : 0
            y: root.down ? Theme.buttonShadowSize : 0
            width: parent.width
            height: parent.height
            color: root.checked
                   ? (root.hovered ? Theme.buttonPrimaryHover : Theme.buttonPrimary)
                   : (root.hovered ? Theme.buttonSecondaryHover : Theme.buttonSecondary)
            border.width: Theme.frameBorder
            border.color: root.checked ? root.accentColor : Theme.buttonBorder
            Behavior on color {
                ColorAnimation { duration: Theme.durationFast }
            }
            Behavior on border.color {
                ColorAnimation { duration: Theme.durationFast }
            }

            Rectangle {
                anchors.left: parent.left
                anchors.right: parent.right
                anchors.top: parent.top
                height: 1
                color: root.checked ? Qt.rgba(1, 0.3, 0.3, 0.35) : Theme.buttonHighlight
                opacity: root.down ? 0 : 1
            }
        }
    }

    ThemedToolTip {
        parent: root
        x: Math.round((root.width - width) / 2)
        y: root.height + Theme.spacingTiny
        visible: root.tooltip !== "" && root.hovered
        text: root.tooltip
    }
}
