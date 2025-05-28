// import React, { useState } from 'react'
// import { Hexagon } from './Hexagon'
// import { Leaderboard } from './Leaderboard'
// interface MapProps {
//   showControls?: boolean
// }
// // Mock data for hexagons
// const mockHexagons = [
//   {
//     id: 1,
//     x: 0,
//     y: 0,
//     owner: 'StrideRunner42',
//     score: 120,
//     color: 'yellow',
//   },
//   {
//     id: 2,
//     x: 1,
//     y: 0,
//     owner: 'WalkMaster',
//     score: 85,
//     color: 'blue',
//   },
//   {
//     id: 3,
//     x: 0,
//     y: 1,
//     owner: 'HikerPro',
//     score: 95,
//     color: 'green',
//   },
//   {
//     id: 4,
//     x: -1,
//     y: 0,
//     owner: null,
//     score: 0,
//     color: 'gray',
//   },
//   {
//     id: 5,
//     x: 0,
//     y: -1,
//     owner: 'JoggerKing',
//     score: 50,
//     color: 'red',
//   },
//   {
//     id: 6,
//     x: 1,
//     y: -1,
//     owner: 'StrideRunner42',
//     score: 75,
//     color: 'yellow',
//   },
//   {
//     id: 7,
//     x: -1,
//     y: 1,
//     owner: 'WalkMaster',
//     score: 60,
//     color: 'blue',
//   },
//   {
//     id: 8,
//     x: 2,
//     y: -1,
//     owner: null,
//     score: 0,
//     color: 'gray',
//   },
//   {
//     id: 9,
//     x: -2,
//     y: 1,
//     owner: 'HikerPro',
//     score: 40,
//     color: 'green',
//   },
//   {
//     id: 10,
//     x: 1,
//     y: 1,
//     owner: 'JoggerKing',
//     score: 30,
//     color: 'red',
//   },
//   {
//     id: 11,
//     x: -1,
//     y: -1,
//     owner: null,
//     score: 0,
//     color: 'gray',
//   },
//   {
//     id: 12,
//     x: 2,
//     y: 0,
//     owner: 'StrideRunner42',
//     score: 110,
//     color: 'yellow',
//   },
// ]
// export const Map: React.FC<MapProps> = ({ showControls = true }) => {
//   const [selectedHexagon, setSelectedHexagon] = useState<number | null>(null)
//   const [showLeaderboard, setShowLeaderboard] = useState(false)
//   const handleHexagonClick = (id: number) => {
//     setSelectedHexagon(id)
//     setShowLeaderboard(true)
//   }
//   const closeLeaderboard = () => {
//     setShowLeaderboard(false)
//     setSelectedHexagon(null)
//   }
//   const getSelectedHexagonData = () => {
//     return mockHexagons.find((hex) => hex.id === selectedHexagon)
//   }
//   // Calculate positions based on hexagon grid geometry
//   const calculatePosition = (x: number, y: number) => {
//     const hexSize = 40
//     const hexWidth = hexSize * 2
//     const hexHeight = Math.sqrt(3) * hexSize
//     const posX = x * (hexWidth * 0.75) + 150
//     const posY = y * hexHeight + (x % 2 === 0 ? 0 : hexHeight / 2) + 150
//     return {
//       posX,
//       posY,
//     }
//   }
//   return (
//     <div className="relative w-full h-full bg-gray-900 overflow-hidden">
//       <div className="absolute inset-0 flex items-center justify-center">
//         <div className="relative w-full h-full">
//           {mockHexagons.map((hex) => {
//             const { posX, posY } = calculatePosition(hex.x, hex.y)
//             return (
//               <Hexagon
//                 key={hex.id}
//                 id={hex.id}
//                 x={posX}
//                 y={posY}
//                 owner={hex.owner}
//                 color={hex.color}
//                 isSelected={selectedHexagon === hex.id}
//                 onClick={handleHexagonClick}
//               />
//             )
//           })}
//           {/* User location indicator */}
//           <div
//             className="absolute w-6 h-6 bg-white rounded-full shadow-lg border-2 border-yellow-500 animate-pulse"
//             style={{
//               left: '150px',
//               top: '150px',
//               transform: 'translate(-50%, -50%)',
//             }}
//           />
//           {/* Current position label */}
//           <div
//             className="absolute px-2 py-1 bg-gray-800 rounded-md text-xs text-white"
//             style={{
//               left: '150px',
//               top: '170px',
//               transform: 'translateX(-50%)',
//             }}
//           >
//             Current Location
//           </div>
//         </div>
//       </div>
//       {showLeaderboard && selectedHexagon && (
//         <Leaderboard
//           hexagon={getSelectedHexagonData()}
//           onClose={closeLeaderboard}
//         />
//       )}
//     </div>
//   )
// }
