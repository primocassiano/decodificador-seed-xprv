package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
	"strings"
	"time"

	// Import the local copy of LND's aezeed package
	"aezeed_tool/internal/aezeed"

	// Need hdkeychain for BIP32 derivation
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
)

// BitcoinGenesisDate is the timestamp of Bitcoin's genesis block.
// Copied from LND's aezeed package for birthday calculation.
var BitcoinGenesisDate = time.Unix(1231006505, 0)

// timeFromBitcoinDaysGenesis computes the time from the number of days since the
// Bitcoin genesis block.
func timeFromBitcoinDaysGenesis(days uint16) time.Time {
	// Convert days to duration
	dayDuration := time.Duration(days) * 24 * time.Hour

	// Add the duration to the genesis date.
	return BitcoinGenesisDate.Add(dayDuration)
}

// ToHDRootKey derives the BIP32 HD root key (master node) from the seed's entropy.
// Takes entropy as [16]byte as returned by LND's aezeed.
func ToHDRootKey(entropy [16]byte, net *chaincfg.Params) (string, error) {
	// BIP32 uses the entropy directly as the seed for the master key.
	// hdkeychain.NewMaster expects a []byte slice.
	masterKey, err := hdkeychain.NewMaster(entropy[:], net)
	 if err != nil {
		 return "", fmt.Errorf("failed to derive master HD key: %w", err)
	 }

	 // Return the extended private key in base58 format.
	 return masterKey.String(), nil
}

// CLI interface for the aezeed tool
func runCLI() {
	var choice int
	for {
		fmt.Println("\nAezeed Tool - Menu Principal")
		fmt.Println("1. Gerar nova seed aezeed")
		fmt.Println("2. Decodificar mnemônico existente")
		fmt.Println("3. Sair")
		fmt.Print("Escolha uma opção (1-3): ")

		_, err := fmt.Scanf("%d\n", &choice) // Read integer and newline
		 if err != nil {
			 fmt.Println("Entrada inválida. Por favor, digite um número (1-3).")
			 // Clear buffer if needed
			 reader := bufio.NewReader(os.Stdin)
			 // Read and discard the rest of the line
			 _, _ = reader.ReadString('\n') // Use single quotes for byte literal
			 continue
		 }

		 switch choice {
		 case 1:
			 generateNewSeed()
		 case 2:
			 decodeMnemonic()
		 case 3:
			 fmt.Println("Saindo...")
			 return
		 default:
			 fmt.Println("Opção inválida. Por favor, tente novamente.")
		 }
	 }
}

// generateNewSeed creates a new seed and displays its mnemonic and HD root key
func generateNewSeed() {
	var passphraseInput string

	fmt.Println("\n--- Gerar Nova Seed ---")
	fmt.Print("Digite uma frase-senha (deixe em branco para usar a padrão \"aezeed\"): ")

	// Read the full line for passphrase
	scanner := bufio.NewScanner(os.Stdin)
	 if scanner.Scan() {
		 passphraseInput = scanner.Text()
	 }

	 // Use nil passphrase if input is empty, LND's package handles the default.
	 var passphrase []byte
	 if passphraseInput != "" {
		 passphrase = []byte(passphraseInput)
	 }

	 // Generate 16 bytes of random entropy as [16]byte array
	 var entropy [16]byte
	 if _, err := rand.Read(entropy[:]); err != nil { // Read into the slice view of the array
		 fmt.Printf("Erro ao gerar entropia: %v\n", err)
		 return
	 }

	 // Get current time for the seed birthday
	 birthTime := time.Now()

	 // Create a new CipherSeed instance using the LND package
	 // aezeed.New expects version, *[16]byte, time.Time
	 // Version 0 is standard for aezeed
	 seed, err := aezeed.New(0, &entropy, birthTime)
	 if err != nil {
		 fmt.Printf("Erro ao criar nova seed (LND): %v\n", err)
		 return
	 }

	 // Convert the seed to mnemonic using the LND package
	 // ToMnemonic returns Mnemonic type which is [NumMnemonicWords]string
	 mnemonicArray, err := seed.ToMnemonic(passphrase)
	 if err != nil {
		 fmt.Printf("Erro ao gerar mnemônico (LND): %v\n", err)
		 return
	 }

	 // Convert the array to a slice for joining
	 mnemonicSlice := mnemonicArray[:]
	 fmt.Printf("\nMnemônico Gerado: %s\n", strings.Join(mnemonicSlice, " "))
	 
	 // Access Entropy field directly
	 fmt.Printf("Entropia (hex): %x\n", seed.Entropy)
	 
	 // Access Birthday field (uint16) and convert to time.Time for display
	 birthdayTime := timeFromBitcoinDaysGenesis(seed.Birthday)
	 fmt.Printf("Timestamp: %s (%d dias ABE)\n", birthdayTime.Format(time.RFC3339), seed.Birthday)

	 // Derive and print the HD root key (Bitcoin Mainnet, Legacy BIP44)
	 rootKey, err := ToHDRootKey(seed.Entropy, &chaincfg.MainNetParams)
	 if err != nil {
		 fmt.Printf("Erro ao derivar chave HD root: %v\n", err)
		 return
	 }
	 fmt.Printf("Chave HD Root (xprv): %s\n", rootKey)
}

// decodeMnemonic decodes an existing mnemonic and displays its details
func decodeMnemonic() {
	var mnemonicInput, passphraseInput string

	fmt.Println("\n--- Decodificar Mnemônico Existente ---")
	fmt.Println("Digite o mnemônico (24 palavras separadas por espaço):")

	// Read the full line for mnemonic
	scanner := bufio.NewScanner(os.Stdin)
	 if scanner.Scan() {
		 mnemonicInput = scanner.Text()
	 }

	 fmt.Print("Digite a frase-senha (deixe em branco para usar a padrão \"aezeed\"): ")
	 if scanner.Scan() {
		 passphraseInput = scanner.Text()
	 }

	 // Use nil passphrase if input is empty, LND's package handles the default.
	 var passphrase []byte
	 if passphraseInput != "" {
		 passphrase = []byte(passphraseInput)
	 }

	 // Split mnemonic string into slice
	 mnemonicWords := strings.Fields(mnemonicInput)
	 if len(mnemonicWords) != aezeed.NumMnemonicWords {
		 fmt.Printf("Erro: Mnemônico deve ter %d palavras, mas tem %d\n", aezeed.NumMnemonicWords, len(mnemonicWords))
		 return
	 }

	 // Convert slice to array for LND's API
	 var mnemonicArray aezeed.Mnemonic
	 for i, word := range mnemonicWords {
		 if i < len(mnemonicArray) {
			 mnemonicArray[i] = word
		 }
	 }

	 // Decode the mnemonic using the LND package
	 // ToCipherSeed is a method on the Mnemonic type
	 decryptedSeed, err := mnemonicArray.ToCipherSeed(passphrase)
	 if err != nil {
		 fmt.Printf("Erro ao descriptografar seed (LND): %v\n", err)
		 return
	 }

	 fmt.Printf("\nEntropia Decodificada (hex): %x\n", decryptedSeed.Entropy)
	 // Access Birthday field (uint16) and convert to time.Time for display
	 birthdayTime := timeFromBitcoinDaysGenesis(decryptedSeed.Birthday)
	 fmt.Printf("Timestamp: %s (%d dias ABE)\n", birthdayTime.Format(time.RFC3339), decryptedSeed.Birthday)

	 // Derive and print the HD root key
	 rootKey, err := ToHDRootKey(decryptedSeed.Entropy, &chaincfg.MainNetParams)
	 if err != nil {
		 fmt.Printf("Erro ao derivar chave HD root: %v\n", err)
		 return
	 }
	 fmt.Printf("Chave HD Root (xprv): %s\n", rootKey)
}

func main() {
	// No need to load wordlist explicitly anymore
	// LND's aezeed package uses an embedded wordlist

	fmt.Println("\nAezeed Tool - Ferramenta para geração e decodificação de seeds aezeed")
	fmt.Println("Esta ferramenta permite gerar novas seeds aezeed e decodificar mnemônicos existentes.")
	fmt.Println("Utiliza a implementação aezeed do LND para compatibilidade com o site https://guggero.github.io/cryptography-toolkit/#!/aezeed")

	// Run the CLI interface
	 runCLI()
}
