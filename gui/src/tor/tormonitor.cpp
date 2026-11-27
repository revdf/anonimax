#include "tor/tormonitor.h"

#include <QFile>
#include <QVariantMap>

namespace {
constexpr quint16 ControlPort = 9051;
constexpr auto    ControlHost = "127.0.0.1";

const QStringList CookiePaths = {
    "/run/tor/control.authcookie",
    "/var/run/tor/control.authcookie",
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
  m_launchedAt.clear();
  m_requestQueue.clear();
  m_reqPending  = false;
  m_inMultiLine = false;
  m_responseLines.clear();
  m_buffer.clear();
  emit changed();
}

QVariantList TorMonitor::circuits() const {
  QVariantList list;
  const QDateTime now = QDateTime::currentDateTime();
  for (const auto &c : m_circuits) {
    // Build display path with country codes where available
    QStringList hops;
    for (int i = 0; i < c.nodes.size(); i++) {
      const QString cc = (i < c.countries.size() && !c.countries[i].isEmpty())
                         ? QStringLiteral("[%1] ").arg(c.countries[i])
                         : QString();
      hops << cc + c.nodes[i];
    }

    // Age
    const QString age = formatAge(c.builtAt);

    // Build time
    QString buildTime;
    if (c.buildMs >= 0)
      buildTime = QStringLiteral("%1s").arg(c.buildMs / 1000.0, 0, 'f', 1);

    QVariantMap m;
    m[QStringLiteral("id")]        = c.id;
    m[QStringLiteral("display")]   = hops.join(QStringLiteral(" → "));
    m[QStringLiteral("age")]       = age;
    m[QStringLiteral("buildTime")] = buildTime;
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

  // Events (650) are always processed, even while waiting for a response.
  if (line.startsWith(QLatin1String("650 "))) {
    const QString event = line.mid(4);
    if (event.startsWith(QLatin1String("BW ")))
      parseBandwidth(event.mid(3));
    else if (event.startsWith(QLatin1String("CIRC ")))
      parseCircuit(event.mid(5));
    return;
  }

  // Handle GETINFO responses while a request is pending.
  if (!m_reqPending)
    return;

  if (m_inMultiLine) {
    if (line == QLatin1String("."))
      m_inMultiLine = false;
    else
      m_responseLines.append(line);
    return;
  }

  if (line.startsWith(QLatin1String("250+"))) {
    // Start of multi-line response — data follows on subsequent lines.
    m_inMultiLine = true;
    return;
  }

  if (line.startsWith(QLatin1String("250-"))) {
    m_responseLines.append(line.mid(4));
    return;
  }

  if (line.startsWith(QLatin1String("250 ")) || line == QLatin1String("250 OK")) {
    auto handler = m_requestQueue.head().handler;
    m_requestQueue.dequeue();
    m_reqPending = false;
    handler(m_responseLines);
    m_responseLines.clear();
    sendNextRequest();
    return;
  }

  // Error response (5xx) — discard and continue.
  if (!m_requestQueue.isEmpty() && line.startsWith(QLatin1String("5"))) {
    m_requestQueue.dequeue();
    m_reqPending = false;
    m_responseLines.clear();
    sendNextRequest();
  }
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

  const int     id    = parts[0].toInt();
  const QString state = parts[1];

  if (state == QLatin1String("LAUNCHED")) {
    m_launchedAt[id] = QDateTime::currentDateTime();

  } else if (state == QLatin1String("BUILT") && parts.size() >= 3) {
    const auto [fps, nodes] = parsePath(parts[2]);
    const int buildMs = m_launchedAt.contains(id)
                        ? static_cast<int>(m_launchedAt[id].msecsTo(QDateTime::currentDateTime()))
                        : -1;
    m_launchedAt.remove(id);
    m_circuits[id] = { id, fps, nodes,
                       QStringList(fps.size(), QString()),
                       QDateTime::currentDateTime(), buildMs };
    lookupCountries(id);

  } else if (state == QLatin1String("CLOSED") ||
             state == QLatin1String("FAILED")) {
    m_circuits.remove(id);
    m_launchedAt.remove(id);
  }
  emit changed();
}

void TorMonitor::lookupCountries(int circId) {
  if (!m_circuits.contains(circId))
    return;
  const QStringList fps = m_circuits[circId].fps;

  for (int i = 0; i < fps.size(); i++) {
    const QString fp  = fps[i];
    const int     idx = i;

    // Step 1: get the relay's network status entry to find its IP.
    enqueueRequest(
        QStringLiteral("GETINFO ns/id/%1").arg(fp),
        [this, circId, idx](QStringList lines) {
          QString ip;
          for (const QString &l : lines) {
            if (l.startsWith(QLatin1String("r "))) {
              const QStringList f = l.split(QLatin1Char(' '));
              if (f.size() >= 7) ip = f[6];
              break;
            }
          }
          if (ip.isEmpty() || !m_circuits.contains(circId))
            return;

          // Step 2: resolve the IP to a country code.
          enqueueRequest(
              QStringLiteral("GETINFO ip-to-country/%1").arg(ip),
              [this, circId, idx](QStringList lines) {
                if (lines.isEmpty() || !m_circuits.contains(circId))
                  return;
                const int eq = lines.first().lastIndexOf(QLatin1Char('='));
                if (eq < 0) return;
                const QString cc = lines.first().mid(eq + 1).trimmed().toUpper();
                if (cc.length() == 2) {
                  while (m_circuits[circId].countries.size() <= idx)
                    m_circuits[circId].countries.append(QString());
                  m_circuits[circId].countries[idx] = cc;
                  emit changed();
                }
              }
          );
        }
    );
  }
}

void TorMonitor::enqueueRequest(const QString &cmd,
                                std::function<void(QStringList)> handler) {
  m_requestQueue.enqueue({ cmd, std::move(handler) });
  if (!m_reqPending)
    sendNextRequest();
}

void TorMonitor::sendNextRequest() {
  if (m_requestQueue.isEmpty())
    return;
  m_reqPending = true;
  m_socket.write((m_requestQueue.head().command + QLatin1String("\r\n")).toUtf8());
}

QPair<QStringList, QStringList> TorMonitor::parsePath(const QString &pathStr) {
  QStringList fps, names;
  for (const QString &hop : pathStr.split(QLatin1Char(','))) {
    // Format: $<40-char-hex>~<name>  or  $<40-char-hex>=<name>
    const int start  = hop.startsWith(QLatin1Char('$')) ? 1 : 0;
    const int tilde  = hop.indexOf(QLatin1Char('~'));
    const int eq     = hop.indexOf(QLatin1Char('='));
    const int sep    = (tilde > 0) ? tilde : eq;
    if (sep > 0) {
      fps   << hop.mid(start, sep - start).toUpper();
      names << hop.mid(sep + 1);
    } else {
      const QString fp = hop.mid(start).toUpper();
      fps   << fp;
      names << fp.left(8).toLower();
    }
  }
  return { fps, names };
}

QString TorMonitor::formatAge(const QDateTime &builtAt) {
  if (!builtAt.isValid()) return {};
  const qint64 secs = builtAt.secsTo(QDateTime::currentDateTime());
  if (secs < 60)  return QStringLiteral("%1s").arg(secs);
  if (secs < 3600) return QStringLiteral("%1m %2s").arg(secs / 60).arg(secs % 60);
  return QStringLiteral("%1h %2m").arg(secs / 3600).arg((secs % 3600) / 60);
}

QString TorMonitor::cookieHex() {
  for (const QString &path : CookiePaths) {
    QFile f(path);
    if (f.open(QIODevice::ReadOnly))
      return QString::fromLatin1(f.readAll().toHex());
  }
  return {};
}
