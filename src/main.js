import App from './App.svelte';
import { Notary, readFileFromBrowser, sha256 } from './notary';
// import { readFileFromBrowser } from './notary';
import {
  createAccount,
  Client,
  encodeHex,
  decodeHex
} from 'orbs-client-sdk/dist/index.es';

const SENDER_PUBLIC_KEY = 'sender_public_key';
const SENDER_PRIVATE_KEY = 'sender_private_key';
const SENDER_ADDRESS = 'sender_address';

if (!localStorage.getItem(SENDER_PUBLIC_KEY)) {
  const sender = createAccount();
  localStorage.setItem(SENDER_PUBLIC_KEY, encodeHex(sender.publicKey));
  localStorage.setItem(SENDER_PRIVATE_KEY, encodeHex(sender.privateKey));
  localStorage.setItem(SENDER_ADDRESS, sender.address);
}

const publicKey = decodeHex(localStorage.getItem(SENDER_PUBLIC_KEY));
const privateKey = decodeHex(localStorage.getItem(SENDER_PRIVATE_KEY));
const orbsClient = new Client(
  process.env.ORBS_NODE_ADDRESS,
  process.env.ORBS_VCHAIN,
  'TEST_NET'
);

const actions = new Notary(orbsClient, 'Notary', publicKey, privateKey, true);

const app = new App({
  target: document.body,
  props: {
    actions,
    readFileFromBrowser,
    sha256
  }
});

export default app;
