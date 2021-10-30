Nginx RTMP Server & Basic Authentication via `rtmpauthbot`
==========================================================

an automated installation of an nginx based rtmp server built from source configured with a basic authentication system and discord notifications.

Developed for deployment on CentOS Stream 8 or recent release of Fedora (linux x86_64)  
Requires Ansible 2.9+  

Install:
```
ansible-playbook -i root@<serverip>, setup.yml
```

