#pragma once

#include <QObject>
#include <QProcess>
#include <QStringList>

class CommandClient final : public QObject {
  Q_OBJECT

public:
  explicit CommandClient(QObject *parent = nullptr);
  ~CommandClient() override;

  void requestStatus();
  void requestIP();
  void runSafeKill();
  void runPanic();
  void runAutochange(bool start, int minutes = 10);

signals:
  void statusReady(bool anonsurfActive, bool torActive, const QString &dnsState,
                   const QStringList &dnsAddresses);
  void statusFailed(const QString &message);
  void ipReady(const QString &address, bool isTor);
  void ipFailed(const QString &message);
  void safeKillFinished(bool success, const QString &message);
  void panicFinished(bool success, const QString &message);
  void autochangeFinished(bool success, bool started, const QString &message);

private:
  void handleStatusFinished(int exitCode, QProcess::ExitStatus exitStatus);
  void handleIPFinished(int exitCode, QProcess::ExitStatus exitStatus);
  void handleSafeKillFinished(int exitCode, QProcess::ExitStatus exitStatus);
  void handlePanicFinished(int exitCode, QProcess::ExitStatus exitStatus);
  void handleAutochangeFinished(int exitCode, QProcess::ExitStatus exitStatus);
  static QString processError(QProcess &process, const QString &fallback);

  QProcess m_statusProcess;
  QProcess m_ipProcess;
  QProcess m_safeKillProcess;
  QProcess m_panicProcess;
  QProcess m_autochangeProcess;
  bool m_autochangeStarting = false;
};
