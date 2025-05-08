# Ferramenta Aezeed

Esta ferramenta é um programa de linha de comando (CLI) escrito em Go para gerar e decodificar seeds de carteira no formato aezeed, compatível com a implementação usada pelo LND (Lightning Network Daemon) e o site [cryptography-toolkit](https://guggero.github.io/cryptography-toolkit/#!/aezeed).

## Funcionalidades

*   **Geração de Novas Seeds:** Cria uma nova seed aezeed (mnemônico de 24 palavras) com entropia aleatória e timestamp atual.
*   **Decodificação de Mnemônicos:** Decodifica um mnemônico aezeed existente (24 palavras) para recuperar a entropia e o timestamp originais.
*   **Derivação de Chave HD:** Para ambas as funcionalidades (geração e decodificação), deriva e exibe a chave HD root (xprv) correspondente no formato base58, seguindo o padrão BIP32/44 para Bitcoin (legacy).
*   **Suporte a Frase-senha:** Permite o uso de uma frase-senha opcional para proteger a seed. Se nenhuma frase-senha for fornecida, a padrão "aezeed" será usada.
*   **Compatibilidade:** Utiliza o código-fonte do pacote `aezeed` do LND para garantir total compatibilidade com a implementação de referência.

## Como Compilar e Executar

**Pré-requisitos:**

*   Go (versão 1.21 ou superior recomendada)
*   Git (para clonar o repositório, se necessário)

**Passos:**

1.  **Obtenha o Código:**
    *   Se você recebeu os arquivos em um arquivo `.zip`, descompacte-o.
    *   Certifique-se de que a estrutura de diretórios seja:
        ```
        aezeed_tool/
        ├── internal/
        │   └── aezeed/
        │       ├── cipherseed.go
        │       ├── errors.go
        │       └── wordlist.go
        ├── go.mod
        ├── go.sum
        └── main.go
        ```

2.  **Navegue até o Diretório:** Abra seu terminal ou prompt de comando e navegue até o diretório `DECODIFICADOR_SEED_XPRV`:
    ```bash
    cd path/to/DECODIFICADOR_SEED_XPRV
    ```

3.  **Compile (Opcional):** Você pode criar um executável independente:
    ```bash
    go build -o DECODIFICADOR_SEED_XPRV_cli
    ```
    Isso criará um arquivo chamado `DECODIFICADOR_SEED_XPRV_cli` (ou `DECODIFICADOR_SEED_XPRV_cli.exe` no Windows) no diretório atual.

4.  **Execute:**
    *   **Usando `go run` (sem compilar):**
        ```bash
        go run main.go
        ```
    *   **Executando o arquivo compilado:**
        ```bash
        ./DECODIFICADOR_SEED_XPRV_cli 
        ```
        (No Windows, use: `DECODIFICADOR_SEED_XPRV_cli.exe`)

## Como Usar

Ao executar o programa, você verá um menu principal:

```
Aezeed Tool - Menu Principal
1. Gerar nova seed aezeed
2. Decodificar mnemônico existente
3. Sair
Escolha uma opção (1-3):
```

*   **Opção 1: Gerar nova seed aezeed**
    *   O programa solicitará uma frase-senha. Você pode digitar uma frase-senha personalizada ou pressionar Enter para usar a padrão ("aezeed").
    *   Ele exibirá o mnemônico de 24 palavras gerado, a entropia em formato hexadecimal, o timestamp de criação e a chave HD root (xprv) correspondente.

*   **Opção 2: Decodificar mnemônico existente**
    *   O programa solicitará que você digite o mnemônico de 24 palavras, separadas por espaço.
    *   Em seguida, solicitará a frase-senha associada a esse mnemônico (ou deixe em branco para a padrão "aezeed").
    *   Se o mnemônico e a frase-senha estiverem corretos, ele exibirá a entropia decodificada em formato hexadecimal, o timestamp de criação e a chave HD root (xprv) correspondente.
    *   **Importante:** Certifique-se de inserir exatamente 24 palavras.

*   **Opção 3: Sair**
    *   Encerra o programa.

## Detalhes da Implementação

Para garantir a máxima compatibilidade com o site de referência e a implementação do LND, esta ferramenta incorpora diretamente o código-fonte relevante do pacote `aezeed` do LND (`cipherseed.go`, `wordlist.go`, `errors.go`). Isso assegura que a lógica de criptografia, derivação de chave scrypt, manipulação de mnemônico e checksum seja idêntica à usada pelo LND.

A derivação da chave HD root utiliza a biblioteca `github.com/btcsuite/btcd/btcutil/hdkeychain`, configurada para a rede principal do Bitcoin (`chaincfg.MainNetParams`) e o formato legacy (BIP32/44), gerando a chave `xprv`.

