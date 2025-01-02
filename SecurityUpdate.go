package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Função principal do programa
func main() {
	// Verifica o sistema operacional em execução
	switch runtime.GOOS {
	case "linux":
		// Configura atualizações automáticas no Linux
		err := configureLinux()
		if err != nil {
			log.Printf("Erro ao configurar atualizações no Linux: %v", err)
		} else {
			fmt.Println("Atualizações de segurança configuradas com sucesso no Linux.")
		}
	case "windows":
		// Configura atualizações automáticas no Windows
		err := configureWindows()
		if err != nil {
			log.Printf("Erro ao configurar atualizações no Windows: %v", err)
		} else {
			fmt.Println("Atualizações de segurança configuradas com sucesso no Windows.")
		}
	default:
		// Informa se o sistema operacional não é suportado
		log.Printf("Sistema operacional não suportado: %s", runtime.GOOS)
	}

	fmt.Println("\nPara fechar a aplicação, pressione Enter.")

	// Aguarda a entrada do usuário para encerrar a aplicação
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}

// Configurações específicas para sistemas Linux
func configureLinux() error {
	// Detecta a distribuição Linux em uso
	distro, err := getLinuxDistro()
	if err != nil {
		return fmt.Errorf("erro ao detectar a distribuição Linux: %v", err)
	}

	fmt.Printf("Distribuição Linux detectada: %s\n", distro)

	// Define ações com base na distribuição detectada
	switch strings.ToLower(distro) {
	case "ubuntu", "debian", "kali":
		// Atualiza sistemas baseados em Debian, Ubuntu e Kali
		err := updateDebianBased(distro)
		if err != nil {
			return fmt.Errorf("erro ao atualizar sistemas baseados em Debian: %v", err)
		}
	default:
		// Informa se a distribuição não é suportada para atualizações automáticas
		return fmt.Errorf("distribuição Linux não suportada para atualizações automáticas de segurança: %s", distro)
	}

	return nil
}

// Atualiza sistemas baseados em Debian, Ubuntu e Kali
func updateDebianBased(distro string) error {
	fmt.Printf("Aplicando atualizações de segurança no sistema baseado em %s...\n", distro)

	// Atualiza a lista de pacotes disponíveis
	fmt.Println("Atualizando a lista de pacotes...")
	cmd := exec.Command("sudo", "apt-get", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("falha ao executar apt-get update: %v", err)
	}

	// Verifica se o pacote unattended-upgrades está instalado
	if !commandExists("unattended-upgrades") {
		fmt.Println("unattended-upgrades não está instalado. Instalando...")
		cmd := exec.Command("sudo", "apt-get", "install", "unattended-upgrades", "-y")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("erro ao instalar unattended-upgrades: %v", err)
		}
	}

	// Configura o arquivo de configuração do unattended-upgrades
	fmt.Println("Configurando o arquivo 50unattended-upgrades...")
	err := configureUnattendedUpgrades(distro)
	if err != nil {
		return fmt.Errorf("erro ao configurar unattended-upgrades: %v", err)
	}

	// Ativa e reinicia o serviço de atualizações automáticas
	fmt.Println("Ativando e reiniciando o serviço de atualizações automáticas...")
	cmd = exec.Command("sudo", "systemctl", "enable", "--now", "unattended-upgrades.service")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("falha ao ativar o serviço unattended-upgrades: %v", err)
	}

	return nil
}

// Configura o arquivo 50unattended-upgrades para aplicar apenas atualizações de segurança
func configureUnattendedUpgrades(string) error {
	unattendedPath := "/etc/apt/apt.conf.d/50unattended-upgrades"

	// Caminho para o backup do arquivo de configuração
	backupPath := unattendedPath + ".bak"

	// Verifica se o backup já existe
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		// Lê o conteúdo do arquivo original
		input, err := os.ReadFile(unattendedPath)
		if err != nil {
			return fmt.Errorf("falha ao ler %s: %v", unattendedPath, err)
		}
		// Cria um backup do arquivo original
		err = os.WriteFile(backupPath, input, 0644)
		if err != nil {
			return fmt.Errorf("falha ao criar backup do arquivo: %v", err)
		}
		fmt.Printf("Backup criado em %s\n", backupPath)
	} else {
		fmt.Printf("Backup já existe em %s\n", backupPath)
	}

	// Define a configuração desejada para o unattended-upgrades
	config := `// Unattended-Upgrade::Origins-Pattern controla quais pacotes são atualizados.
//
// Atualizações de segurança para Debian e Kali
Unattended-Upgrade::Origins-Pattern {
		"o=Debian,codename=${distro_codename},label=Debian";
		"o=Debian,codename=${distro_codename}-security,label=Debian-Security";
		"o=Kali,codename=${distro_codename},label=Kali-Security";
};

// Pacotes a serem excluídos de atualizações automáticas
Unattended-Upgrade::Package-Blacklist {
// Exemplo de pacotes que você pode querer evitar atualizar automaticamente
// "linux-";
// "libc6";
// Adicione os nomes dos pacotes aqui, se necessário
};

// Ativar a correção automática de problemas com dpkg
Unattended-Upgrade::AutoFixInterruptedDpkg "true";

// Fazer atualizações mínimas, dividindo em pequenos passos
Unattended-Upgrade::MinimalSteps "true";

// Remover pacotes não usados automaticamente após a atualização
Unattended-Upgrade::Remove-New-Unused-Dependencies "true";

// Não reiniciar automaticamente após a atualização
Unattended-Upgrade::Automatic-Reboot "false";

// Limitar a largura de banda para downloads (exemplo: 70 kbps)
//Acquire::http::Dl-Limit "70";

// Habilitar logging detalhado
Unattended-Upgrade::Verbose "true";

// Habilitar modo de depuração para diagnósticos
Unattended-Upgrade::Debug "true";
`

	// Escreve a configuração no arquivo 50unattended-upgrades
	err := os.WriteFile(unattendedPath, []byte(config), 0644)
	if err != nil {
		return fmt.Errorf("falha ao escrever no arquivo %s: %v", unattendedPath, err)
	}
	fmt.Printf("Arquivo %s configurado.\n", unattendedPath)

	return nil
}

// Verifica se um comando está disponível no sistema
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// Detecta a distribuição Linux em uso lendo o arquivo /etc/os-release
func getLinuxDistro() (string, error) {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "", fmt.Errorf("não foi possível ler /etc/os-release: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		// Procura pela linha que começa com "ID="
		if strings.HasPrefix(line, "ID=") {
			id := strings.TrimPrefix(line, "ID=")
			id = strings.Trim(id, `"`) // Remove aspas, se houver
			return id, nil
		}
	}

	// Verifica erros durante a leitura do arquivo
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("erro ao escanear /etc/os-release: %v", err)
	}

	return "", fmt.Errorf("não foi possível detectar a distribuição Linux")
}

// Configurações específicas para sistemas Windows
func configureWindows() error {
	// Informa que a configuração de atualizações automáticas no Windows está iniciando
	fmt.Println("Configurando atualizações automáticas de segurança no Windows...")

	// Script PowerShell para configurar atualizações automáticas
	script := `
# Verificar se o módulo PSWindowsUpdate está instalado
if (-not (Get-Module -ListAvailable -Name PSWindowsUpdate)) {
    try {
        # Instalar o provedor de pacote NuGet não interativamente
        Install-PackageProvider -Name NuGet -MinimumVersion 2.8.5.201 -Force -Scope AllUsers

        # Instalar o módulo PSWindowsUpdate não interativamente
        Install-Module -Name PSWindowsUpdate -Force -Scope AllUsers -ErrorAction Stop
    } catch {
        Write-Error "Falha ao instalar o módulo PSWindowsUpdate ou o provedor de NuGet: $_"
        exit 1
    }
}

# Importar o módulo PSWindowsUpdate para uso
Import-Module PSWindowsUpdate

# Opcional: Configurar políticas de atualização se aplicável
# Por exemplo, definir atualizações automáticas via registry ou GPO

# Definir o nome da tarefa agendada para atualizações automáticas
$taskName = "AutomaticSecurityUpdates"

# Verificar se a tarefa agendada já existe
if (-not (Get-ScheduledTask -TaskName $taskName -ErrorAction SilentlyContinue)) {
    try {
        # Definir a ação a ser executada pela tarefa agendada
        $action = New-ScheduledTaskAction -Execute 'PowerShell.exe' -Argument '-NoProfile -WindowStyle Hidden -Command "Import-Module PSWindowsUpdate; Get-WindowsUpdate -Install -AcceptAll"'
        
        # Definir o gatilho (trigger) para a tarefa agendada (diariamente às 11h)
        $trigger = New-ScheduledTaskTrigger -Daily -At 11am
        
        # Definir o principal (usuário e nível de execução) para a tarefa agendada
        $principal = New-ScheduledTaskPrincipal -UserId "SYSTEM" -RunLevel Highest
        
        # Registrar a tarefa agendada com as ações, gatilho e principal definidos
        Register-ScheduledTask -TaskName $taskName -Action $action -Trigger $trigger -Principal $principal -ErrorAction Stop
        Write-Output "Tarefa agendada '$taskName' criada com sucesso."
    } catch {
        Write-Error "Falha ao criar a tarefa agendada: $_"
        exit 1
    }
} else {
    Write-Output "A tarefa agendada '$taskName' já existe. Nenhuma ação necessária."
}
`

	// Cria um arquivo temporário para armazenar o script PowerShell
	tempFile, err := os.CreateTemp("", "update_script_*.ps1")
	if err != nil {
		return fmt.Errorf("falha ao criar arquivo temporário: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Garante que o arquivo temporário seja removido após a execução

	// Escreve o script PowerShell no arquivo temporário
	_, err = tempFile.Write([]byte(script))
	if err != nil {
		return fmt.Errorf("falha ao escrever no arquivo temporário: %v", err)
	}
	tempFile.Close()

	// Executa o script PowerShell com permissões adequadas
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", tempFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("falha ao executar o script PowerShell: %v", err)
	}

	return nil
}
