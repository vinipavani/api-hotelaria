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

- `make setup` - Realiza a primeira carga das imagens, instala bibliotecas Go e aplica automaticamente as migrações iniciais de tabelas + seed.
- `make seed` - Popula o banco de dados com uma carga realista de 5 hotéis e 100 quartos via Factory.
- `make run` - Sobe os containers da API e do Banco exibindo os logs em tempo real.
- `make test` - Inicializa uma instância paralela (`db_test`), ergue as tabelas via script Go nativo e executa toda a suíte de testes unitários e de integração de forma 100% isolada.
- `make db-migrate` - Varre o diretório virtual do Go (`embed`) e aplica novas tabelas ou colunas criadas no banco de desenvolvimento.
- `make db-rollback` - Desfaz o histórico de alterações do banco. Aceita o parâmetro dinâmico `steps` (Ex: `make db-rollback steps=3`).
- `make create-migration` - Gera um par sequencial de arquivos `.up.sql` e `.down.sql` na pasta de migrações (Ex: `make create-migration name=add_document`).
- `make stop` - Pausa os containers liberando a memória RAM do Linux, mas mantém todos os dados salvos.
- `make clean` - Remove permanentemente os containers e apaga os dados/volumes do PostgreSQL (útil para resets).

---

### 🗺️ Diagrama Entidade-Relacionamento (DER)

<img width="1216" height="331" alt="image" src="https://github.com/user-attachments/assets/653f1e4e-57e4-4ca5-b485-1ec86ab71bb9" />

### 🔗 Regras de Relacionamento e Cardinalidade

A arquitetura do banco de dados foi desenhada para refletir com precisão os processos do negócio hoteleiro, utilizando restrições de chaves estrangeiras rígidas para impedir dados órfãos e garantir a consistência total do sistema.

#### 1. Hotéis (`hotels`) ➡️ Quartos (`rooms`)
* **Entidade de Primeira Classe:** A tabela de hotéis é um recurso independente e livre. Ela não depende de nenhuma outra tabela para nascer no sistema.
* **Pertencimento Estrito (1:N):** Cada quarto obrigatoriamente **precisa pertencer a um único hotel** (`hotel_id NOT NULL`). Não existe quarto órfão ou compartilhado entre estabelecimentos.
* **Cardinalidade:** Um hotel pode nascer totalmente zerado e ter **0 ou N quartos** associados a ele ao longo do tempo.
* **Deleção em Cascata:** Se um hotel for removido do sistema, o PostgreSQL aciona nativamente o comportamento `ON DELETE CASCADE`. Isso significa que todos os quartos atrelados a ele serão sumariamente deletados em background, blindando o banco contra lixo e inconsistências.
* **Unicidade Local Controlada:** Embora dois hotéis diferentes possam ter um quarto com o número `"0101"`, o sistema proíbe que o mesmo hotel tenha dois quartos iguais. Isso é travado pela restrição de unicidade composta (`unique_hotel_room_number`), que amarra o ID do hotel com o número textual do quarto.

#### 2. Quartos (`rooms`) ➡️ Reservas/Hospedagens (`bookings`)
* **Dependência Operacional (1:N):** A tabela de reservas é um recurso totalmente dependente. Uma hospedagem não faz sentido sozinha e precisa obrigatoriamente **apontar para um quarto físico existente** (`room_id NOT NULL`).
* **Cardinalidade Histórica:** Um quarto individual pode possuir um histórico de **0 ou N registros de hospedagem** no decorrer dos anos.
* **Regra de Ouro (Trava Dinâmica de Ocupação):** Embora o relacionamento técnico seja de 1 para N, o motor de regras de negócio em Go (Service) junto ao banco impõe uma restrição temporal rigorosa: **um quarto só pode possuir, no máximo, 1 reserva com o status `'em_estadia'` por vez**. Novos registros de entrada ficam bloqueados até que o registro ativo sofra um evento de check-out, mudando seu estado para `'finalizada'`.
* **Ciclo de Vida Amarrado:** Seguindo o fluxo de cascata, caso um quarto específico seja excluído da base (por exemplo, para uma reforma estrutural definitiva), todo o histórico de bookings e estadias daquele quarto também será limpo automaticamente via `ON DELETE CASCADE`.

## 🚀 Endpoints da API (Exemplos de Uso)

Abaixo estão listados todos os endpoints disponíveis na API de Hotelaria, organizados pelo padrão de recursos RESTful, contendo exemplos práticos de requisições utilizando o `curl` e as respectivas respostas esperadas.

---

### 🏨 Domínio de Hotéis

#### 1. Criar um Novo Hotel
* **Rota:** `POST /hotels`
* **Descrição:** Cadastra um estabelecimento principal no sistema.
* **Exemplo de Requisição:**
```bash
curl -i -X POST http://localhost:8080/hotels \
  -H "Content-Type: application/json" \
  -d '{"name": "Copacabana Palace", "city": "Rio de Janeiro"}'
```
* **Exemplo de Resposta (201 Created):**
```json
{
  "id": 1,
  "name": "Copacabana Palace",
  "city": "Rio de Janeiro",
  "created_at": "2026-07-12T15:00:00Z"
}
```

#### 2. Listar Todos os Hotéis
* **Rota:** `GET /hotels`
* **Descrição:** Retorna a coleção completa de hotéis cadastrados no banco de dados.
* **Exemplo de Requisição:**
```bash
curl -i -X GET http://localhost:8080/hotels
```
* **Exemplo de Resposta (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Copacabana Palace",
    "city": "Rio de Janeiro",
    "created_at": "2026-07-12T15:00:00Z"
  }
]
```

---

### 🔑 Domínio de Quartos

#### 3. Cadastrar um Quarto para um Hotel
* **Rota:** `POST /hotels/:id/rooms`
* **Descrição:** Cadastra um novo quarto associado ao ID do hotel informado na URL. O número do quarto é gerado e formatado sequencialmente de forma atômica pelo banco de dados (ex: `"0001"`, `"0002"`).
* **Exemplo de Requisição:**
```bash
curl -i -X POST http://localhost:8080/hotels/1/rooms \
  -H "Content-Type: application/json" \
  -d '{"type": "suite", "capacity": 4, "per_night_value": 350.00}'
```
* **Exemplo de Resposta (201 Created):**
```json
{
  "id": 1,
  "hotel_id": 1,
  "number": "0001",
  "type": "suite",
  "capacity": 4,
  "per_night_value": 350,
  "created_at": "2026-07-12T15:05:00Z"
}
```

#### 4. Listar Quartos de um Hotel
* **Rota:** `GET /hotels/:id/rooms`
* **Descrição:** Lista os quartos de um hotel. Permite passar o Query Parameter `disponivel=true` para omitir quartos ocupados em tempo real.
* **Exemplo de Requisição:**
```bash
curl -i -X GET http://localhost:8080/hotels/1/rooms
```
* **Exemplo de Resposta (200 OK):**
```json
[
  {
    "id": 1,
    "hotel_id": 1,
    "number": "0002",
    "type": "double",
    "capacity": 2,
    "per_night_value": 200,
    "created_at": "2026-07-12T15:10:00Z"
  }
]
```

---

### 📋 Domínio de Reservas e Hospedagens

#### 5. Efetuar Check-in (Entrada do Hóspede)
* **Rota:** `POST /rooms/:id/check-in`
* **Descrição:** Abre uma nova hospedagem ativa com status `'em_estadia'` vinculada ao quarto solicitado.
* **Exemplo de Requisição:**
```bash
curl -i -X POST http://localhost:8080/rooms/1/check-in \
  -H "Content-Type: application/json" \
  -d '{
    "guest_name": "Vinicius Dev",
    "guest_document": "123.456.789-00",
    "check_in": "2026-07-12"
  }'
```
* **Exemplo de Resposta (201 Created):**
```json
{
  "id": 1,
  "room_id": 1,
  "guest_name": "Vinicius Dev",
  "guest_document": "123.456.789-00",
  "status": "em_estadia",
  "check_in": "2026-07-12",
  "created_at": "2026-07-12T15:20:00Z"
}
```

#### 6. Efetuar Check-out (Saída do Hóspede)
* **Rota:** `POST /rooms/:id/check-out`
* **Descrição:** Encerra a estadia ativa do quarto informado, atualizando o status para `'finalizada'` e registrando a data de saída.
* **Exemplo de Requisição:**
```bash
curl -i -X POST http://localhost:8080/rooms/1/check-out \
  -H "Content-Type: application/json" \
  -d '{"check_out": "2026-07-15"}'
```
* **Exemplo de Resposta (200 OK):**
```json
{
  "id": 1,
  "room_id": 1,
  "guest_name": "Vinicius Dev",
  "guest_document": "123.456.789-00",
  "status": "finalizada",
  "check_in": "2026-07-12",
  "check_out": "2026-07-15",
  "created_at": "2026-07-12T15:20:00Z"
}
```

## 🚀 Arquitetura e Diferenciais
- **Domain-Driven Directory Structure:** Organização de pastas limpa e isolada por domínios de negócio inspirada em DDD.
- **Injeção de Dependência:** Gerenciamento seguro de conexões de banco e configurações de rede sem uso de variáveis globais desordenadas.
- **Ambiente de Alta Produtividade:** Hot-reload automático com a ferramenta **Air** (salve o código e o Docker recompila em milissegundos).
- **Variáveis Tipadas:** Configurações centralizadas e validadas logo na inicialização da aplicação através do `godotenv`.
