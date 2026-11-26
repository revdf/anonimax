import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

Dialog {
    id: root

    required property var config
    readonly property int maxWidth: 480

    parent: Overlay.overlay
    anchors.centerIn: parent
    width: parent ? Math.min(root.maxWidth, parent.width - Theme.contentMargin * 2)
                  : root.maxWidth
    modal: true
    focus: true
    padding: Theme.spacingLarge
    standardButtons: Dialog.NoButton
    header: null

    onOpened: {
        pathField.text = root.config.vpnConfigPath || ""
        errorText.text = ""
    }

    contentItem: ColumnLayout {
        spacing: Theme.spacingMedium

        StatusBadge {
            Layout.alignment: Qt.AlignLeft
            text: "VPN"
            accentColor: Theme.accent
        }

        Text {
            text: "Provide the path to your VPN configuration file (.ovpn for OpenVPN, "
                  + ".conf for WireGuard). The VPN will connect before Tor starts."
            color: Theme.textSecondary
            font.family: Theme.uiFont
            font.pixelSize: Theme.bodySize
            wrapMode: Text.WordWrap
            Layout.fillWidth: true
        }

        Rectangle {
            Layout.fillWidth: true
            height: 1
            color: Theme.borderLight
        }

        GridLayout {
            Layout.fillWidth: true
            columns: 2
            columnSpacing: Theme.spacingMedium
            rowSpacing: Theme.spacingSmall

            Text {
                text: "Config file"
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
            }

            TextField {
                id: pathField
                Layout.fillWidth: true
                placeholderText: "/etc/vpn/myvpn.ovpn"
                color: Theme.textPrimary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
                background: Rectangle {
                    color: Theme.buttonSecondary
                    border.width: 1
                    border.color: pathField.activeFocus ? Theme.accent : Theme.buttonBorder
                }
            }
        }

        Text {
            Layout.fillWidth: true
            text: "Supported: OpenVPN (.ovpn), WireGuard (.conf)"
            color: Theme.textSecondary
            font.family: Theme.uiFont
            font.pixelSize: Theme.detailSize
            font.italic: true
        }

        Text {
            id: errorText
            Layout.fillWidth: true
            text: ""
            color: Theme.danger
            font.family: Theme.uiFont
            font.pixelSize: Theme.detailSize
            wrapMode: Text.WordWrap
            visible: text !== ""
        }

        RowLayout {
            Layout.fillWidth: true
            Layout.topMargin: Theme.spacingSmall
            spacing: Theme.spacingSmall

            Item { Layout.fillWidth: true }

            ActionButton {
                Layout.preferredWidth: 120
                text: "Cancel"
                iconName: "close"
                onClicked: root.reject()
            }

            ActionButton {
                Layout.preferredWidth: 120
                text: "Save"
                iconName: "check"
                primary: true
                onClicked: {
                    const p = pathField.text.trim()
                    if (p === "") {
                        errorText.text = "Config file path cannot be empty."
                        return
                    }
                    root.config.setSetting("vpn_config_path", p)
                    root.config.setSetting("use_vpn", "true")
                    root.accept()
                }
            }
        }
    }

    background: Rectangle {
        color: Theme.bgPanel
        border.width: Theme.frameBorder
        border.color: Theme.accent
    }
}
