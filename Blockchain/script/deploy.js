const hre = require("hardhat");

async function main() {
  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying contracts with the account:", deployer.address);

  const PrescriptionNFT = await hre.ethers.getContractFactory("PrescriptionNFT");
  const prescriptionNFT = await PrescriptionNFT.deploy(deployer.address);
  
  // Wait for the deployment transaction to be mined
  await prescriptionNFT.waitForDeployment();
  
  // Get the deployed contract address
  const deployedAddress = await prescriptionNFT.getAddress();
  console.log("PrescriptionNFT deployed to:", deployedAddress);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });