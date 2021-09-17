# supervisor

## 

http://supervisord.org/installing.html

## 

```
sudo pip install --upgrade pip
sudo pip install -r requirements.txt
```

## supervisord 

```
http://supervisord.org/running.html#supervisorctl-actions

supervisord  supervisord  supervisorctl ：

supervisord  Supervisord  。
supervisorctl stop programxxx (programxxx)，programxxx  [program:beepkg]   beepkg。
supervisorctl start programxxx 
supervisorctl restart programxxx 
supervisorctl stop groupworker:   groupworker (start,restart )
supervisorctl stop all  ：start、restart、stop 。
supervisorctl reload   。
supervisorctl update   。

   stop   reload  update 。
```

## relayd

```
echo_supervisord_conf > supervisord.conf
sudo mv echo_supervisord_conf > /etc/supervisord.conf
sudo mkdir /etc/supervisord.conf.d
sudo echo "[include]" >> /etc/supervisord.conf
sudo echo "files = /etc/supervisord.conf.d/*.conf" >> /etc/supervisord.conf
sudo cp supervisord_relayd.conf /etc/supervisord.conf.d/
```

