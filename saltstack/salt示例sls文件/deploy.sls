deploy:
  file.directory:
    - name: /data/{{ pillar['service'] }}/{{ pillar['service'] }}.{{ pillar['version'] }}
    - makedirs: True
    - user: {{ pillar['service'] }}
    - group: {{ pillar['service'] }}


pull:
  file.managed:
    - name: /data/{{ pillar['service'] }}/{{ pillar['service'] }}.{{ pillar['version'] }}/{{ pillar['service'] }}.tar.gz
    - source: salt://pre/{{ pillar['service'] }}.tar.gz
    - user: {{ pillar['service'] }}
    - group: {{ pillar['service'] }}


tar:
  cmd.run:
    - name: cd /data/{{ pillar['service'] }}/{{ pillar['service'] }}.{{ pillar['version'] }} && tar xf {{ pillar['service'] }}.tar.gz


link:
  file.symlink:
    - name: /data/{{ pillar['service'] }}/{{ pillar['service'] }}
    - target: /data/{{ pillar['service'] }}/{{ pillar['service'] }}.{{ pillar['version'] }}
    - force: True
delete:
  cmd.run:
    - name: cd /data/{{ pillar['service'] }} && find ./ -type d -name "{{ pillar['service'] }}.*"|xargs du -h --time --max-depth=0|head -n -5|awk '{print $4}'|xargs rm -rf
    # 删除前5个
    # - name: cd /data/{{ pillar['service'] }} && rm -rf {{ pillar['service'] }}.`expr {{ pillar['version'] }} - 5`

restart:
  cmd.run:
    - name: /etc/init.d/{{ pillar['service'] }} restart

