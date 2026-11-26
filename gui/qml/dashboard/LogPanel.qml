import QtQuick
import QtQuick.Layouts
import QtQuick.Controls

Rectangle {
    id: root

    required property var logModel

    color: Theme.bgPanel
    border.width: 1
    border.color: Theme.borderLight

    ColumnLayout {
        anchors.fill: parent
        anchors.margins: Theme.spacingMedium
        spacing: Theme.spacingSmall

        RowLayout {
            Layout.fillWidth: true

            Text {
                text: "ACTIVITY LOG"
                color: Theme.textSecondary
                font.family: Theme.uiFont
                font.pixelSize: Theme.badgeSize
                font.weight: Font.Medium
            }

            Item { Layout.fillWidth: true }

            Text {
                text: "CLEAR"
                color: Theme.accent
                font.family: Theme.uiFont
                font.pixelSize: Theme.badgeSize
                font.bold: true

                MouseArea {
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: root.logModel.clear()
                }
            }
        }

        Rectangle {
            Layout.fillWidth: true
            height: 1
            color: Theme.borderLight
        }

        ListView {
            id: logView

            Layout.fillWidth: true
            Layout.fillHeight: true
            clip: true
            model: root.logModel
            spacing: 1
            verticalLayoutDirection: ListView.BottomToTop

            delegate: RowLayout {
                width: logView.width
                spacing: Theme.spacingSmall

                Text {
                    text: model.time
                    color: Theme.textSecondary
                    font.family: Theme.uiFont
                    font.pixelSize: Theme.detailSize
                }

                Text {
                    text: model.message
                    color: model.entryType === "error"   ? Theme.danger
                         : model.entryType === "success" ? Theme.success
                         : model.entryType === "warn"    ? Theme.warning
                         : Theme.textPrimary
                    font.family: Theme.uiFont
                    font.pixelSize: Theme.detailSize
                    Layout.fillWidth: true
                    wrapMode: Text.WordWrap
                }
            }
        }
    }
}
