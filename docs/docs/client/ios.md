---
sidebar_position: 2
---

# iOS

Using Meemaw in your iOS app? Let's get your started with our iOS SDK!

:::info
Meemaw's iOS SDK works in concert with the [Argent Labs Web3 Library](https://github.com/argentlabs/web3.swift). The library is well maintained and well used. The library also allows us to offer a very elegant solution as it lets us define a custom wallet type. However, that makes our SDK more coupled than we like. If you have an opinion on what other Web3 Swift library we should support, please reach out on Discord!
:::

## Install SDK in your Xcode project

Adding our iOS SDK to your app is as simple as adding a Package Dependency using our iOS SDK repo: `https://github.com/getmeemaw/meemaw-ios`

You can find more info about adding packages in [Xcode docs](https://developer.apple.com/documentation/xcode/adding-package-dependencies-to-your-app).

## Use SDK in Swift

### Init Meemaw & Web3 then get wallet

Your first step is to initialise Meemaw and to "get" a wallet. If no wallet currently exists for a given user, the SDK will generate a new one in concert with the server (TSS process). If one exists, the SDK will use it.

```swift
import meemaw
import web3

guard let clientUrl = URL(string: "JSON-RPC_API_URL") else { return }
let client = EthereumHttpClient(url: clientUrl, network: EthereumNetwork.sepolia) // using the Sepolia test net in this case

// Create or recover a wallet for this user
let meemaw = Meemaw(server: "meemaw-url")
let wallet = try await meemaw.GetWallet(auth: TOKEN) // will recover the wallet if exists for the user or create a new one
```

The first two lines import Meemaw and Argent Web3 libraries. The second two lines initialise the web3 library.

The last two lines initialise the library with the server address and "get" the wallet. One interesting bit is the `TOKEN`. It represents the user connexion and depends on your Auth mechanism. For example, if you're using Supabase, the token is the Supabase's `access_token`. Behind the scenes, Meemaw will authenticate the user using that token and only procede if the user exists AND is logged in.

Once you create a wallet for a user, potentially at sign up, you will most probably want to display or store the public key for that wallet. You can easily do that by calling :

```javascript
let publicKey = wallet.From()
```

### Sign and send normal transaction (Ethereum)

In order to send funds, the first step is to craft the raw transaction:

```swift
// Craft the transaction (adapt to  your needs) 
let transaction = EthereumTransaction(
    from: wallet.From(),
    to: "RECIPIENT_ADDRESS",
    value: 10000000000000,
    data: Data(),
    nonce: 1,
    gasPrice: 10,
    gasLimit: 21000,
    chainId: 1
)
```

then to sign it using the Meemaw wallet and send it:

```swift
// Sign and send the transaction with wallet
let txHash = try await client.eth_sendRawTransaction(transaction, withAccount: wallet)
print("Look at my fabulous transaction: \(txHash)")
```

The SDK will automatically sign the transaction in concert with the server (TSS process) and provide the signed transaction to the Web3 library sending it through. Note that we need to provide the `chainId` in the transaction, in order to avoid replay attacks on other chains. 

Also, you should get the `nonce` and `gasPrice` using the web3 library.

### Sign and send smart contract call (Ethereum)

The way Meemaw's iOS SDK is integrated with the [Argent Labs Web3 library](https://github.com/argentlabs/web3.swift), you can basically use that library normally, just with Meemaw's wallet. The example above for a normal transaction shows you the way. For smart contracts, follow [their readme here](https://github.com/argentlabs/web3.swift?tab=readme-ov-file#smart-contracts-static-types).

### Sign message (Ethereum)

Our iOS SDK provides a way to sign Ethereum messages following the [eip-191](https://eips.ethereum.org/EIPS/eip-191) and [eip-712](https://eips.ethereum.org/EIPS/eip-712) standards, i.e. the Ethereum prefix will be automatically added and the message will be properly hashed:

```swift
let signature = try wallet.SignEthMessage(message: Data)
```

### Sign bytes (all ECDSA blockchains)

Our iOS SDK also provides a way to sign arbitrary bytes:

```swift
let signature = try wallet.SignBytes(message: Data)
```

### Multi-device

Before you start using it, it's probably important you learn [how multi-device works](/docs/multi-device).

Multi-device works in two steps: 
1. starting registration on a new device
2. accepting on an existing device.

#### Register device

You actually don't have to do anything. The function we used above (`meemaw.GetWallet(...)`) automatically recognizes that it's a new device when a wallet already exists.

However, we recommend that you provide two callback functions: one called when the function starts the registration process, and one called when the process is finished. This allows you to display a call to action in the UI, prompting the user to use his other device to confirm.

Here is how you would adapt the code above:

```swift
import meemaw
import web3

let meemaw = Meemaw(server: "meemaw-url")
let wallet = try await meemaw.GetWallet(
    auth: TOKEN, 
    callbackRegisterStarted: {
        // RegisterDevice started => prompt user for confirmation on existing device
    },
    callbackRegisterDone: {
        // RegisterDevice done => hide prompt
    })
```

This will recover the wallet if it already exists on the device, create a new one if the user doesn't have one at all, or start the multi-device process if the user already has a wallet created on another device.

#### Accept device

On the existing device, it is now time to confirm the new device with a simple call:

```javascript
await wallet.AcceptDevice()
```

That's it, now the multi-device process will happen with both devices and the server communicating with each other. At the end of the process, the new device will have it's own part of the overall MPC wallet, ready for operations.

### Export private key

You can offer your users to export their private key:

```swift
let privateKey = try wallet.Export()
```

This is useful if your users want to manage their wallet outside of your platform, for example.

:::warning
Meemaw wallets are MPC wallets. Until you decide otherwise, **no private key exists**, transactions are signed through a collaboration process between the server and the client.

However, **it is possible to generate the private key corresponding to the wallet**. The way it works is the following: the client sends its TSS share to the server, which then combines the client share with its own share to form a private key.

It is important to understand that by doing so, **you completely loose the security and decentralization benefits of MPC wallets**: anyone with access to this private key can spend the funds contained in the wallet.

It is also important to note that **the MPC wallet still exist**, it's just that the private key now exists as well and can bypass the whole TSS process.
:::