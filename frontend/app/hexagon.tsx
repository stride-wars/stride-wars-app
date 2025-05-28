// import React from 'react'
// interface HexagonProps {
//   id: number
//   x: number
//   y: number
//   owner: string | null
//   color: string
//   isSelected: boolean
//   onClick: (id: number) => void
// }
// export const Hexagon: React.FC<HexagonProps> = ({
//   id,
//   x,
//   y,
//   owner,
//   color,
//   isSelected,
//   onClick,
// }) => {
//   // Size of the hexagon
//   const size = 40
//   // Calculate points for hexagon shape
//   const getHexPoints = () => {
//     const points = []
//     for (let i = 0; i < 6; i++) {
//       const angle = (i * 60 * Math.PI) / 180
//       const pointX = size * Math.cos(angle)
//       const pointY = size * Math.sin(angle)
//       points.push(`${pointX},${pointY}`)
//     }
//     return points.join(' ')
//   }
//   // Determine fill color based on owner
//   const getFillColor = () => {
//     if (owner === null) return '#4B5563' // Gray for unclaimed
//     switch (color) {
//       case 'yellow':
//         return '#F59E0B'
//       case 'blue':
//         return '#3B82F6'
//       case 'green':
//         return '#10B981'
//       case 'red':
//         return '#EF4444'
//       default:
//         return '#4B5563'
//     }
//   }
//   return (
//     <div
//       className="absolute cursor-pointer transform -translate-x-1/2 -translate-y-1/2"
//       style={{
//         left: `${x}px`,
//         top: `${y}px`,
//       }}
//       onClick={() => onClick(id)}
//     >
//       <svg
//         width={size * 2}
//         height={size * 2}
//         viewBox={`${-size} ${-size} ${size * 2} ${size * 2}`}
//       >
//         <polygon
//           points={getHexPoints()}
//           fill={getFillColor()}
//           stroke={isSelected ? '#FFFFFF' : '#1F2937'}
//           strokeWidth={isSelected ? 3 : 2}
//           opacity={owner === null ? 0.5 : 0.8}
//         />
//         {owner && (
//           <text
//             x="0"
//             y="0"
//             textAnchor="middle"
//             dominantBaseline="middle"
//             fill="white"
//             fontSize="10"
//             fontWeight="bold"
//           >
//             {owner.substring(0, 2)}
//           </text>
//         )}
//       </svg>
//     </div>
//   )
// }
