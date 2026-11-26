pragma Singleton
import QtQuick

QtObject {
    readonly property color accent: "#B01010"
    readonly property color accentDim: "#8B0000"
    readonly property color accentBright: "#FF3B30"
    readonly property color warning: "#F4C95D"
    readonly property color danger: "#FF453A"
    readonly property color info: "#B01010"
    readonly property color neonCyan: "#FF6B6B"
    readonly property color bgWindow: "#120404"
    readonly property color bgPanel: "#1A0606"
    readonly property color bgPanelAlt: "#230808"
    readonly property color bgBottom: "#0A0202"
    readonly property color textPrimary: "#FCF4F4"
    readonly property color textSecondary: "#C4A0A0"
    readonly property color textOnAccent: "#FFF7F7"
    readonly property color borderLight: Qt.rgba(1, 0.22, 0.22, 0.18)
    readonly property color buttonText: "#FFCCCC"
    readonly property color buttonPrimary: "#B01010"
    readonly property color buttonPrimaryHover: "#990000"
    readonly property color buttonSecondary: "#1A0606"
    readonly property color buttonSecondaryHover: "#2D0A0A"
    readonly property color buttonBorder: "#CC1515"
    readonly property color buttonShadow: "#3D0000"
    readonly property color buttonSecondaryShadow: "#0C0202"
    readonly property color buttonHighlight: Qt.rgba(1, 0.3, 0.3, 0.28)
    readonly property color scanline: Qt.rgba(0, 0, 0, 0.34)
    readonly property color vignette: Qt.rgba(0, 0, 0, 0.58)
    readonly property color grid: Qt.rgba(0.65, 0.1, 0.1, 0.09)

    readonly property string uiFont: "Inter"
    readonly property int titleSize: 26
    readonly property int headingSize: 18
    readonly property int bodySize: 15
    readonly property int buttonSize: 13
    readonly property int badgeSize: 12
    readonly property int detailSize: 13
    readonly property int spacingTiny: 5
    readonly property int spacingSmall: 8
    readonly property int spacingMedium: 10
    readonly property int spacingLarge: 16
    readonly property int contentMargin: 16
    readonly property int frameBorder: 2
    readonly property int buttonHeight: 40
    readonly property int buttonShadowSize: 3
    readonly property int iconButtonSize: 38
    readonly property int badgeHeight: 24
    readonly property int badgePaddingX: 9
    readonly property int durationFast: 150
    readonly property int durationMedium: 220
    readonly property int scanlineStep: 4
}
