import QtQuick

Canvas {
    id: root

    property var  downHistory: []
    property var  upHistory:   []
    property color downColor: Theme.success
    property color upColor:   Theme.accent

    onDownHistoryChanged: requestPaint()
    onUpHistoryChanged:   requestPaint()

    function formatRate(bytes) {
        if (bytes >= 1048576) return (bytes / 1048576).toFixed(1) + " MB/s"
        if (bytes >= 1024)    return (bytes / 1024).toFixed(1)    + " KB/s"
        return bytes.toFixed(0) + " B/s"
    }

    onPaint: {
        var ctx = getContext("2d")
        ctx.clearRect(0, 0, width, height)

        const all   = downHistory.concat(upHistory)
        const maxV  = Math.max(1024, ...all)
        const pad   = 4
        const w     = width  - pad * 2
        const h     = height - pad * 2

        // Grid lines
        ctx.strokeStyle = Qt.rgba(1, 0.2, 0.2, 0.12)
        ctx.lineWidth   = 1
        for (let g = 1; g <= 3; g++) {
            const gy = pad + (h / 4) * g
            ctx.beginPath()
            ctx.moveTo(pad, gy)
            ctx.lineTo(pad + w, gy)
            ctx.stroke()
        }

        // Draw a smooth line for a data series
        function drawLine(data, color) {
            if (!data || data.length < 2) return
            const step = w / (data.length - 1)
            ctx.strokeStyle = color
            ctx.lineWidth   = 2
            ctx.shadowColor = color
            ctx.shadowBlur  = 4
            ctx.beginPath()
            ctx.moveTo(pad, pad + h - (data[0] / maxV) * h)
            for (let i = 1; i < data.length; i++) {
                const x0  = pad + (i - 1) * step
                const y0  = pad + h - (data[i - 1] / maxV) * h
                const x1  = pad + i * step
                const y1  = pad + h - (data[i] / maxV) * h
                const cpx = (x0 + x1) / 2
                ctx.bezierCurveTo(cpx, y0, cpx, y1, x1, y1)
            }
            ctx.stroke()
            ctx.shadowBlur = 0

            // Fill area under line
            ctx.lineTo(pad + (data.length - 1) * step, pad + h)
            ctx.lineTo(pad, pad + h)
            ctx.closePath()
            const grad = ctx.createLinearGradient(0, pad, 0, pad + h)
            grad.addColorStop(0, Qt.rgba(
                parseInt(color.toString().slice(1,3), 16) / 255,
                parseInt(color.toString().slice(3,5), 16) / 255,
                parseInt(color.toString().slice(5,7), 16) / 255,
                0.18))
            grad.addColorStop(1, Qt.rgba(0, 0, 0, 0))
            ctx.fillStyle = grad
            ctx.fill()
        }

        drawLine(upHistory,   upColor)
        drawLine(downHistory, downColor)
    }
}
