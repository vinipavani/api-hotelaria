# 🏨 API de Hotelaria
Uma API REST moderna, escalável e robusta para gerenciamento de hotelaria, desenvolvida em **Go** (Golang) com o framework **Gin** e banco de dados **PostgreSQL**, totalmente containerizada com **Docker**.

---

## 🛠️ Pré-requisitos Técnicos
Para rodar este projeto localmente, a sua máquina Linux precisa ter instalado apenas:
1. [Docker](https://docker.com) (Versão estável atualizada)
2. [Docker Compose](https://docker.com) (Integrado ao Docker)
3. **Make** (Utilitário nativo do Linux para automação)

> 💡 *Nota: Você **não** precisa instalar o compilador do Go ou o banco PostgreSQL diretamente no seu sistema operacional físico. O Docker cuida de todo o ecossistema.*

---

## 🏃‍♂️ Passo a Passo para Instalação e Execução

Siga as etapas abaixo no terminal do seu Linux para subir o projeto do zero:

### 1. Clonar o Repositório e Acessar a Pasta
```bash
git clone https://github.com/vinipavani/api-hotelaria
cd api-hotelaria
```

### 2. Configurar as Variáveis de Ambiente
Crie o arquivo `.env` na raiz do projeto:
```bash
cp .env.example .env
```
Abra o arquivo `.env` gerado e certifique-se de que as credenciais do banco estejam configuradas para o ambiente Docker:
```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@db:5432/api_db?sslmode=disable
```

### 3. Executar o Setup Inicial (Apenas na primeira vez)
O comando abaixo baixará as imagens oficiais, criará a estrutura do módulo Go e instalará todas as dependências (`Gin`, `pgxpool`, `godotenv`) automaticamente dentro do container:
```bash
make setup
```

### 4. Inicializar a API
Com o setup concluído, sempre que for programar no projeto, basta dar a partida oficial:
```bash
make run
```
*A API estará ativa em `http://localhost:8080` e monitorando qualquer alteração que você fizer no código.*

---

## 🎛️ Painel de Controle (Comandos Disponíveis)
O projeto utiliza um `Makefile` para resumir operações longas do Docker em comandos simples:

- `make setup` - Realiza a primeira carga das imagens, instala bibliotecas Go e cria o ambiente.
- `make run` - Sobe os containers da API e do Banco exibindo os logs em tempo real.
- `make stop` - Pausa os containers liberando a memória RAM do Linux, mas mantém todos os dados salvos.
- `make clean` - Remove permanentemente os containers e apaga os dados/volumes do PostgreSQL (útil para resets).

## 🚀 Arquitetura e Diferenciais
- **Domain-Driven Directory Structure:** Organização de pastas limpa e isolada por domínios de negócio inspirada em DDD.
- **Injeção de Dependência:** Gerenciamento seguro de conexões de banco e configurações de rede sem uso de variáveis globais desordenadas.
- **Ambiente de Alta Produtividade:** Hot-reload automático com a ferramenta **Air** (salve o código e o Docker recompila em milissegundos).
- **Variáveis Tipadas:** Configurações centralizadas e validadas logo na inicialização da aplicação através do `godotenv`.