import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

Dialog {
    id: root

    required property var config
    readonly property int maxWidth: 480
    readonly property bool isSocks5: typeSelector.currentIndex === 0
    readonly property bool isHttps:  typeSelector.currentIndex === 2

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
        // Pre-fill with existing values
        if (root.config.socks5Address !== "") {
            typeSelector.currentIndex = 0
            var parts = root.config.socks5Address.split(":")
            hostField.text = parts[0] || ""
            portField.text = parts.length > 1 ? parts[1] : "1080"
        } else if (root.config.httpAddress !== "") {
            typeSelector.currentIndex = 1
            var hparts = root.config.httpAddress.split(":")
            hostField.text = hparts[0] || ""
            portField.text = hparts.length > 1 ? hparts[1] : "8080"
        } else {
            typeSelector.currentIndex = 0
            hostField.text = ""
            portField.text = "1080"
        }
        userField.text = ""
        passField.text = ""
    }

    contentItem: ColumnLayout {
        spacing: Theme.spacingMedium

        StatusBadge {
            Layout.alignment: Qt.AlignLeft
            text: "PROXY"
            accentColor: Theme.accent
        }

        Text {
            text: "Configure the proxy used before connecting to Tor."
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
                text: "Type"
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
            }

            ComboBox {
                id: typeSelector
                Layout.fillWidth: true
                model: ["SOCKS5", "HTTP", "HTTPS"]
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
                onCurrentIndexChanged: portField.text = currentIndex === 0 ? "1080" : "8080"

                contentItem: RowLayout {
                    spacing: 0

                    Text {
                        Layout.fillWidth: true
                        leftPadding: Theme.spacingSmall
                        text: typeSelector.displayText
                        color: Theme.textPrimary
                        font: typeSelector.font
                        verticalAlignment: Text.AlignVCenter
                        elide: Text.ElideRight
                    }

                    Text {
                        text: "▼"
                        color: Theme.textSecondary
                        font.pixelSize: 10
                        rightPadding: Theme.spacingSmall
                        verticalAlignment: Text.AlignVCenter
                    }
                }

                background: Rectangle {
                    color: typeSelector.hovered ? Theme.buttonSecondaryHover : Theme.buttonSecondary
                    border.width: 1
                    border.color: typeSelector.pressed ? Theme.accent : Theme.buttonBorder
                }

                delegate: Rectangle {
                    id: comboDelegate

                    property bool isHovered: false
                    readonly property bool isSelected: typeSelector.currentIndex === index

                    width: typeSelector.width
                    height: Theme.buttonHeight
                    color: isSelected ? Theme.accent
                         : isHovered  ? Theme.buttonSecondaryHover
                         :              Theme.bgPanelAlt

                    Text {
                        anchors.verticalCenter: parent.verticalCenter
                        leftPadding: Theme.spacingSmall
                        text: modelData
                        color: comboDelegate.isSelected ? Theme.textOnAccent : Theme.textPrimary
                        font.family: Theme.uiFont
                        font.pixelSize: Theme.detailSize
                    }

                    MouseArea {
                        anchors.fill: parent
                        hoverEnabled: true
                        cursorShape: Qt.PointingHandCursor
                        onEntered: comboDelegate.isHovered = true
                        onExited:  comboDelegate.isHovered = false
                        onClicked: {
                            typeSelector.currentIndex = index
                            typeSelector.popup.close()
                        }
                    }
                }

                popup: Popup {
                    y: typeSelector.height
                    width: typeSelector.width
                    padding: 0

                    contentItem: ListView {
                        implicitHeight: contentHeight
                        model: typeSelector.popup.visible ? typeSelector.delegateModel : null
                        clip: true
                    }

                    background: Rectangle {
                        color: Theme.bgPanelAlt
                        border.width: 1
                        border.color: Theme.buttonBorder
                    }
                }
            }

            Text {
                text: "Address"
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
            }

            TextField {
                id: hostField
                Layout.fillWidth: true
                placeholderText: "127.0.0.1"
                color: Theme.textPrimary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
                background: Rectangle {
                    color: Theme.buttonSecondary
                    border.width: 1
                    border.color: hostField.activeFocus ? Theme.accent : Theme.buttonBorder
                }
            }

            Text {
                text: "Port"
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
            }

            TextField {
                id: portField
                Layout.fillWidth: true
                placeholderText: root.isSocks5 ? "1080" : "8080"
                validator: IntValidator { bottom: 1; top: 65535 }
                color: Theme.textPrimary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
                background: Rectangle {
                    color: Theme.buttonSecondary
                    border.width: 1
                    border.color: portField.activeFocus ? Theme.accent : Theme.buttonBorder
                }
            }

            Text {
                text: "User (optional)"
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
            }

            TextField {
                id: userField
                Layout.fillWidth: true
                placeholderText: "username"
                color: Theme.textPrimary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
                background: Rectangle {
                    color: Theme.buttonSecondary
                    border.width: 1
                    border.color: userField.activeFocus ? Theme.accent : Theme.buttonBorder
                }
            }

            Text {
                text: "Password (optional)"
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
            }

            TextField {
                id: passField
                Layout.fillWidth: true
                placeholderText: "password"
                echoMode: TextInput.Password
                color: Theme.textPrimary
                font.family: Theme.uiFont
                font.pixelSize: Theme.detailSize
                background: Rectangle {
                    color: Theme.buttonSecondary
                    border.width: 1
                    border.color: passField.activeFocus ? Theme.accent : Theme.buttonBorder
                }
            }

        }

        // Error message
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
                    if (hostField.text.trim() === "") {
                        errorText.text = "Address cannot be empty."
                        return
                    }
                    if (portField.text.trim() === "") {
                        errorText.text = "Port cannot be empty."
                        return
                    }
                    const addr = hostField.text.trim() + ":" + portField.text.trim()
                    const cred = userField.text.trim() !== ""
                                 ? userField.text.trim() + ":" + passField.text
                                 : ""
                    if (root.isSocks5) {
                        root.config.setSetting("socks5_proxy_address", addr)
                        root.config.setSetting("socks5_proxy_credential", cred)
                        root.config.setSetting("http_proxy_address", "")
                    } else {
                        root.config.setSetting("http_proxy_address", addr)
                        root.config.setSetting("http_proxy_credential", cred)
                        root.config.setSetting("proxy_is_https",
                                               root.isHttps ? "true" : "false")
                        root.config.setSetting("socks5_proxy_address", "")
                    }
                    root.config.setSetting("use_proxy", "true")
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
