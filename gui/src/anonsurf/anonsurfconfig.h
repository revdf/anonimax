#pragma once

#include <QObject>
#include <QProcess>
#include <QString>

class AnonSurfConfig final : public QObject {
  Q_OBJECT
  Q_PROPERTY(bool changeMac        READ changeMac        NOTIFY changed)
  Q_PROPERTY(bool renewIP          READ renewIP          NOTIFY changed)
  Q_PROPERTY(bool forceUTC         READ forceUTC         NOTIFY changed)
  Q_PROPERTY(bool blockWebRTC      READ blockWebRTC      NOTIFY changed)
  Q_PROPERTY(bool blockTelemetry   READ blockTelemetry   NOTIFY changed)
  Q_PROPERTY(bool disableCameraMic READ disableCameraMic NOTIFY changed)
  Q_PROPERTY(bool cleanOnStop      READ cleanOnStop      NOTIFY changed)
  Q_PROPERTY(bool useVPN           READ useVPN           NOTIFY changed)
  Q_PROPERTY(bool useProxy         READ useProxy         NOTIFY changed)
  Q_PROPERTY(QString vpnConfigPath READ vpnConfigPath    NOTIFY changed)
  Q_PROPERTY(QString socks5Address READ socks5Address    NOTIFY changed)
  Q_PROPERTY(QString httpAddress   READ httpAddress      NOTIFY changed)
  Q_PROPERTY(bool busy             READ busy             NOTIFY changed)

public:
  explicit AnonSurfConfig(QObject *parent = nullptr);

  bool    changeMac()        const { return m_changeMac; }
  bool    renewIP()          const { return m_renewIP; }
  bool    forceUTC()         const { return m_forceUTC; }
  bool    blockWebRTC()      const { return m_blockWebRTC; }
  bool    blockTelemetry()   const { return m_blockTelemetry; }
  bool    disableCameraMic() const { return m_disableCameraMic; }
  bool    cleanOnStop()      const { return m_cleanOnStop; }
  bool    useVPN()           const { return m_useVPN; }
  bool    useProxy()         const { return m_useProxy; }
  QString vpnConfigPath()    const { return m_vpnConfigPath; }
  QString socks5Address()    const { return m_socks5Address; }
  QString httpAddress()      const { return m_httpAddress; }
  bool    busy()             const { return m_busy; }

  Q_INVOKABLE void refresh();
  Q_INVOKABLE void setSetting(const QString &key, const QString &value);

signals:
  void changed();

private:
  void parseOutput(const QString &output);

  bool    m_changeMac        = true;
  bool    m_renewIP          = true;
  bool    m_forceUTC         = true;
  bool    m_blockWebRTC      = true;
  bool    m_blockTelemetry   = false;
  bool    m_disableCameraMic = false;
  bool    m_cleanOnStop      = false;
  bool    m_useVPN           = false;
  bool    m_useProxy         = false;
  QString m_vpnConfigPath;
  QString m_socks5Address;
  QString m_httpAddress;
  bool    m_busy = false;

  QProcess m_writeProcess;
};
