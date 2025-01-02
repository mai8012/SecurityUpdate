Configurador de Atualizações Automáticas de Segurança
Este projeto em Go automatiza a configuração de atualizações de segurança 
em sistemas operacionais Linux (distribuições baseadas em Debian como Ubuntu, 
Debian e Kali) e Windows. Ele detecta o sistema operacional no qual está sendo 
executado e aplica as configurações necessárias para garantir que as atualizações 
de segurança sejam instaladas automaticamente, melhorando a segurança e a manutenção do sistema.

Funcionalidades
Detecção Automática do Sistema Operacional: Identifica se o sistema é Linux ou Windows e, 
no caso do Linux, detecta a distribuição específica.
Configuração de Atualizações Automáticas no Linux:
Atualiza a lista de pacotes.
Instala e configura o unattended-upgrades para aplicar atualizações de segurança automaticamente.
Habilita e inicia o serviço de atualizações automáticas.
Configuração de Atualizações Automáticas no Windows:
Instala e configura o módulo PowerShell PSWindowsUpdate.
Configura o Windows Update para baixar e instalar automaticamente as atualizações de segurança.
Cria uma tarefa agendada para verificar e instalar atualizações diariamente.

Configuração do Unattended Upgrades/Sistemas Linux
O arquivo de configuração /etc/apt/apt.conf.d/50unattended-upgrades é ajustado 
para garantir que apenas atualizações de segurança sejam aplicadas automaticamente, 
com opções adicionais como remoção de dependências não utilizadas e logging detalhado.

Como compilar:
![image](https://github.com/user-attachments/assets/a5c66e1b-fd3a-41c8-b555-ba52fea54771)
