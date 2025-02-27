const PrescriptionNFT = artifacts.require("PrescriptionNFT");

module.exports = function (deployer, network, accounts) {
  const initialOwner = accounts[0]; // Use the first account as the initial owner
  deployer.deploy(PrescriptionNFT, initialOwner);
};