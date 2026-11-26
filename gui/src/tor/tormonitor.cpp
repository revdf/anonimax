#include "tor/tormonitor.h"

#include <QFile>
#include <QVariantMap>

namespace {
constexpr quint16 ControlPort = 9051;
constexpr auto    ControlHost = "127.0.0.1";

const QStringList CookiePaths = {
    "/run/tor/control.authcookie",
    "/var/run/tor/control.authcookie",
    "/run/tor/control",
};
} // namespace

TorMonitor::TorMonitor(QObject *parent) : QObject(parent) {
  connect(&m_socket, &QTcpSocket::connected,    this, &TorMonitor::onConnected);
  connect(&m_socket, &QTcpSocket::readyRead,    this, &TorMonitor::onReadyRead);
  connect(&m_socket, &QAbstractSocket::errorOccurred,
          this, &TorMonitor::onError);
}

void TorMonitor::start() {
  if (m_state != State::Disconnected)
    return;
  m_status = "Connecting...";
  emit changed();
  m_socket.connectToHost(ControlHost, ControlPort);
}

void TorMonitor::stop() {
  m_socket.disconnectFromHost();
  m_state     = State::Disconnected;
  m_status    = "Not connected";
  m_downRate  = 0;
  m_upRate    = 0;
  m_circuits.clear();
  m_buffer.clear();
  emit changed();
}

QVariantList TorMonitor::circuits() const {
  QVariantList list;
  for (const auto &c : m_circuits) {
    QVariantMap m;
    m[QStringLiteral("id")]    = c.id;
    m[QStringLiteral("nodes")] = c.nodes;
    m[QStringLiteral("path")]  = c.nodes.join(QStringLiteral(" → "));
    list.append(m);
  }
  return list;
}

void TorMonitor::onConnected() {
  m_state = State::Authenticating;
  authenticate();
}

void TorMonitor::authenticate() {
  const QString hex = cookieHex();
  if (hex.isEmpty())
    m_socket.write("AUTHENTICATE \"\"\r\n");
  else
    m_socket.write(QStringLiteral("AUTHENTICATE %1\r\n").arg(hex).toUtf8());
}

void TorMonitor::onReadyRead() {
  m_buffer += QString::fromUtf8(m_socket.readAll());
  int pos;
  while ((pos = m_buffer.indexOf(QLatin1String("\r\n"))) >= 0) {
    const QString line = m_buffer.left(pos);
    m_buffer.remove(0, pos + 2);
    processLine(line);
  }
}

void TorMonitor::onError(QAbstractSocket::SocketError) {
  m_state  = State::Disconnected;
  m_status = "Cannot connect to Tor control port";
  m_buffer.clear();
  emit changed();
}

void TorMonitor::processLine(const QString &line) {
  if (m_state == State::Authenticating) {
    if (line.startsWith(QLatin1String("250"))) {
      m_state = State::Subscribing;
      m_socket.write("SETEVENTS BW CIRC\r\n");
    } else {
      m_state  = State::Disconnected;
      m_status = "Tor control auth failed";
      m_socket.disconnectFromHost();
      emit changed();
    }
    return;
  }

  if (m_state == State::Subscribing) {
    if (line.startsWith(QLatin1String("250"))) {
      m_state  = State::Active;
      m_status = "Connected";
      emit changed();
    }
    return;
  }

  if (!line.startsWith(QLatin1String("650 ")))
    return;

  const QString event = line.mid(4);
  if (event.startsWith(QLatin1String("BW ")))
    parseBandwidth(event.mid(3));
  else if (event.startsWith(QLatin1String("CIRC ")))
    parseCircuit(event.mid(5));
}

void TorMonitor::parseBandwidth(const QString &data) {
  const QStringList parts = data.split(QLatin1Char(' '));
  if (parts.size() < 2)
    return;
  m_downRate = parts[0].toDouble();
  m_upRate   = parts[1].toDouble();

  m_downHistory.append(m_downRate);
  m_upHistory.append(m_upRate);
  if (m_downHistory.size() > HistorySize) m_downHistory.removeFirst();
  if (m_upHistory.size()   > HistorySize) m_upHistory.removeFirst();

  emit changed();
}

void TorMonitor::parseCircuit(const QString &data) {
  const QStringList parts = data.split(QLatin1Char(' '));
  if (parts.size() < 2)
    return;

  const int    id     = parts[0].toInt();
  const QString state = parts[1];

  if (state == QLatin1String("BUILT") && parts.size() >= 3) {
    m_circuits[id] = { id, parsePath(parts[2]) };
  } else if (state == QLatin1String("CLOSED") ||
             state == QLatin1String("FAILED")) {
    m_circuits.remove(id);
  }
  emit changed();
}

QStringList TorMonitor::parsePath(const QString &pathStr) {
  QStringList names;
  for (const QString &hop : pathStr.split(QLatin1Char(','))) {
    const int tilde = hop.indexOf(QLatin1Char('~'));
    const int eq    = hop.indexOf(QLatin1Char('='));
    const int pos   = (tilde > 0) ? tilde : eq;
    if (pos > 0)
      names << hop.mid(pos + 1);
    else if (hop.startsWith(QLatin1Char('$')) && hop.length() > 9)
      names << hop.mid(1, 8);
    else
      names << hop;
  }
  return names;
}

QString TorMonitor::cookieHex() {
  for (const QString &path : CookiePaths) {
    QFile f(path);
    if (f.open(QIODevice::ReadOnly)) {
      return QString::fromLatin1(f.readAll().toHex());
    }
  }
  return {};
}
