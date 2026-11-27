#pragma once

#include <QDateTime>
#include <QMap>
#include <QObject>
#include <QQueue>
#include <QStringList>
#include <QTcpSocket>
#include <QVariantList>
#include <functional>

class TorMonitor final : public QObject {
  Q_OBJECT
  Q_PROPERTY(bool connected    READ connected    NOTIFY changed)
  Q_PROPERTY(QString status    READ status       NOTIFY changed)
  Q_PROPERTY(double downRate   READ downRate     NOTIFY changed)
  Q_PROPERTY(double upRate     READ upRate       NOTIFY changed)
  Q_PROPERTY(QVariantList downHistory READ downHistory NOTIFY changed)
  Q_PROPERTY(QVariantList upHistory   READ upHistory   NOTIFY changed)
  Q_PROPERTY(int circuitCount  READ circuitCount NOTIFY changed)
  Q_PROPERTY(QVariantList circuits READ circuits NOTIFY changed)

public:
  explicit TorMonitor(QObject *parent = nullptr);

  bool         connected()    const { return m_state == State::Active; }
  QString      status()       const { return m_status; }
  double       downRate()     const { return m_downRate; }
  double       upRate()       const { return m_upRate; }
  QVariantList downHistory()  const { return m_downHistory; }
  QVariantList upHistory()    const { return m_upHistory; }
  int          circuitCount() const { return m_circuits.size(); }
  QVariantList circuits()     const;

  Q_INVOKABLE void start();
  Q_INVOKABLE void stop();

signals:
  void changed();

private:
  enum class State { Disconnected, Authenticating, Subscribing, Active };

  struct Circuit {
    int         id;
    QStringList fps;        // hex fingerprints
    QStringList nodes;      // relay display names
    QStringList countries;  // 2-letter codes, "" while pending
    QDateTime   builtAt;
    int         buildMs = -1;
  };

  struct Request {
    QString                               command;
    std::function<void(QStringList)>      handler;
  };

  void onConnected();
  void onReadyRead();
  void onError(QAbstractSocket::SocketError err);
  void processLine(const QString &line);
  void parseBandwidth(const QString &data);
  void parseCircuit(const QString &data);
  void authenticate();
  void enqueueRequest(const QString &cmd, std::function<void(QStringList)> handler);
  void sendNextRequest();
  void lookupCountries(int circId);
  static QPair<QStringList, QStringList> parsePath(const QString &pathStr);
  static QString cookieHex();
  static QString formatAge(const QDateTime &builtAt);

  QTcpSocket   m_socket;
  QString      m_buffer;
  State        m_state  = State::Disconnected;
  QString      m_status = "Not connected";
  double       m_downRate = 0;
  double       m_upRate   = 0;
  QVariantList m_downHistory;
  QVariantList m_upHistory;

  QMap<int, Circuit>   m_circuits;
  QMap<int, QDateTime> m_launchedAt;

  QQueue<Request> m_requestQueue;
  bool            m_reqPending   = false;
  bool            m_inMultiLine  = false;
  QStringList     m_responseLines;

  static constexpr int HistorySize = 60;
};
