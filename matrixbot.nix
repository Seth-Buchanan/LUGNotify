{ lib, pkgs, ... }:

{
systemd.timers."matrix-bot" = {
  wantedBy = [ "timers.target" ];
    timerConfig = {
      OnCalendar = "Wed *-*-* 7:30:00";
      Unit = "matrix-bot.service";
    };
};

systemd.services."matrix-bot" = {
  script = ''
    /home/johns/source/matrixbot/matrixbot --config /etc/matrixbot/config.json &>> /etc/matrixbot/errors.log
  '';
  serviceConfig = {
    Type = "oneshot";
    User = "root";
  };
};



}
